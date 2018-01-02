/**
 * Created by Jingle Chen on 2017/12/7.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import * as _ from 'lodash'
import { Table, Button, Row, Col, Card, Input, Icon, message } from 'antd';
import { fetchData, receiveData } from '../../action';
import { getPros } from '../../axios';
import BreadcrumbCustom from '../../components/BreadcrumbCustom';
import EditableCell from '../../components/cells/EditableCell';



class UserManager extends React.Component {
    state = {
        selectedRowKeys: [],  // Check here to configure the default column
        loading: false,
        usersData: [],
        companiesData: [],
        selectedCompany: '',
        selectedCompanyId: '',
        editable: false,
    };
    componentDidMount() {
        this.start();
    }
    start = () => {
        this.setState({ loading: true });
        this.fetchData();
        this.fetchCompanyList();
    };

    fetchData = () => {
        const { fetchData } = this.props
        let tempTownId
        fetchData({funcName: 'fetchUsers', stateName: 'usersData'}).then(res => {
            if(res === undefined || res.data === undefined || res.data.users === undefined) return
            this.setState({
                usersData: [...res.data.users.map(val => {
                    val.key = val.id;
                    return val;
                })],
                loading: false,
            });
        });
    }

    fetchCompanyList = () => {
        const { fetchData } = this.props
        fetchData({funcName: 'fetchCompanies', stateName: 'companiesData', params: {}}).then(res => {
            if(res === undefined || res.data === undefined || res.data.companies === undefined) return
            this.setState({
                companiesData: [...res.data.companies.map(val => {
                    val.key = val.id;
                    return val;
                })],
                loading: false,
            });
        });
    }


    onNewRowChange = (dataIndex, value) => {
        this.setState({
            [dataIndex]: value,
        })
    }


    onSelectChange = (selectedRowKeys) => {
        if (selectedRowKeys.length > 0) {
            selectedRowKeys = [selectedRowKeys[selectedRowKeys.length - 1]]
        }

        this.setState({ selectedRowKeys });
    };

    onRowClick = (record, index, event) => {
        const { selectedRowKeys, editable } = this.state
        if(record.id === -1 || editable){
            return
        }
        this.setState({
            selectedRowKeys: selectedRowKeys.length > 0 && selectedRowKeys[0] === record.id ? [] : [record.id],
        });
    }

    handleAdd = () => {
        this.setState({
            usersData: [{
                key: -1,
                id: -1,
                name: '',
                phone: '',
                job: '',
                company_id: '',
                //wx_openid: null,
            }, ...this.state.usersData]
        });
    }

    handleCancelEditRow = () => {
        let tmpUsersData = [...this.state.usersData.filter(item => item.id !== -1)]
        this.setState({
            editable: false,
            usersData: tmpUsersData,
            selectedRowKeys: [],
        })
    }

    handleModify = () => {
        this.setState({
            editable: true,
        })
    }

    handleDelete = () => {
        const { fetchData } = this.props
        const { selectedRowKeys, currentPage } = this.state
        if (selectedRowKeys.length === 0) return
        fetchData({
            funcName: 'deleteUser', params: { id: selectedRowKeys[0] }, stateName: 'deleteUserStatus'
            }).then(res => {
                message.success('删除成功')
                this.fetchData()
            }).catch(err => {
                let errRes = err.response
                if(errRes.data && errRes.data.status === 'error'){
                    message.error(errRes.data.error)
                }
            });
    }

    onNewRowChange = (dataIndex, value) => {
        this.setState({
            [dataIndex]: value,
        })
    }

    onRowSave = () => {
        const { editable } = this.state
        if(editable){
            this.onUpdateRowSave()
        }else{
            this.onNewRowSave()
        }
    }

    onUpdateRowSave = () => {
        const { fetchData } = this.props
        const { selectedRowKeys } = this.state
        const keys = _.keys(this.state)
        const PREFIX = 'user.'
        const PREFIX_LEN = PREFIX.length;
        let obj = {}
        for (let key of keys) {
            if(_.startsWith(key, PREFIX)){
                let field = key.substring(PREFIX_LEN)
                obj[field] = this.state[key]
            }
        }
        obj.id = selectedRowKeys[0]

        fetchData({funcName: 'updateUser', params: obj, stateName: 'updateUserStatus'})
            .then(res => {
                message.success('更新成功')
                this.fetchData()
                this.setState({
                    editable: false,
                })
            }).catch(err => {
                let errRes = err.response
                if(errRes.data && errRes.data.status === 'error'){
                    message.error(errRes.data.error)
                }
            });
    }

    onNewRowSave = () => {
        const { fetchData } = this.props
        const keys = _.keys(this.state)
        const PREFIX = 'user.'
        const PREFIX_LEN = PREFIX.length;
        let obj = {}
        for (let key of keys) {
            if(_.startsWith(key, PREFIX)){
                let field = key.substring(PREFIX_LEN)
                obj[field] = this.state[key]
            }
        }

        fetchData({funcName: 'newUser', params: obj, stateName: 'newUseryStatus'})
            .then(res => {
                message.success('创建成功')
                this.fetchData()
                this.setState({
                    editable: false,
                })
            }).catch(err => {
                let errRes = err.response
                if(errRes.data && errRes.data.status === 'error'){
                    message.error(errRes.data.error)
                }
            });
    }

    render() {

        const { loading, selectedRowKeys, companiesData,
            usersData, editable } = this.state;
        const rowSelection = {
            selectedRowKeys,
            onChange: this.onSelectChange,
            type: 'radio',
        };

        let options = [];
        if(companiesData){
            options = [...companiesData.map(item => {item.key = item.id; return item})]
        }

        const hasSelected = selectedRowKeys.length > 0 && selectedRowKeys[0] !== -1

        const userColumns = [{
            title: '用户名',
            dataIndex: 'name',
            width: '15%',
            render: (text, record) => {
                if (record.id === -1 || (editable && record.id === selectedRowKeys[0])) {
                    return <EditableCell dataIndex='user.name' value={record.name} onChange={this.onNewRowChange} />
                }
                return <a>{text}</a>
            }
        }, {
            title: '手机号',
            dataIndex: 'phone',
            width: '20%',
            render: (text, record) => {
                if (record.id === -1 || (editable && record.id === selectedRowKeys[0])) {
                    return <EditableCell dataIndex='user.phone' value={record.phone} onChange={this.onNewRowChange} />
                }
                return <a>{text}</a>
            }
        }, {
            title: '职位',
            dataIndex: 'job',
            width: '15%',
            render: (text, record) => {
                if (record.id === -1 || (editable && record.id === selectedRowKeys[0])) {
                    return <EditableCell dataIndex='user.job'  value={record.job} onChange={this.onNewRowChange} />
                }
                return <a>{text}</a>
            }
        }, {
            title: '所在公司',
            dataIndex: 'company_name',
            width: '20%',
            render: (text, record) => {
                if (record.id === -1 || (editable && record.id === selectedRowKeys[0])) {
                    return <EditableCell dataIndex='user.company_id' value={record.company_id} onChange={this.onNewRowChange}
                        editType="select" valueType="int" options={options} placeholder="请选择公司"/>
                }
                return <a>{text}</a>
            }
        },{
            title: '微信号',
            dataIndex: 'wx_openid',
            width: '30%',
            render: (text, record) => {
                if (record.id === -1 || (editable && record.id === selectedRowKeys[0])) {
                    return <EditableCell type="opt" onSave={this.onRowSave} onCancel={this.handleCancelEditRow}/>
                }
                return <a>{text}</a>
            }
        }];

        
        return (
            <div className="gutter-example">
                <BreadcrumbCustom first="安监管理" second="用户管理" />
                <Row gutter={16}>
                    <Col className="gutter-row" md={24}>
                        <div className="gutter-box">
                            <Card title="用户列表" bordered={false}>
                                <div style={{ marginBottom: 16 }}>
                                    <Button type="primary" onClick={this.handleAdd}
                                            disabled={loading} 
                                    >新增</Button>
                                    <Button type="primary" onClick={this.handleModify}
                                        disabled={!hasSelected}
                                    >修改</Button>
                                    <Button type="primary" onClick={this.handleDelete}
                                        disabled={!hasSelected}
                                    >删除</Button>
                                </div>
                                <Table rowSelection={rowSelection} columns={userColumns} dataSource={usersData}
                                    onRow={(record) => ({
                                        onClick: () => this.onRowClick(record),
                                    })}
                                />
                            </Card>
                        </div>
                    </Col>
                </Row>
            </div>
        );
    }
}

const mapStateToProps = state => {
    const { 
        companiesData = {data: {count:0, companies:[]}}, 
        townsData = {data: {count:0, towns:[]}}, 
        countries = {data: {count:0, countries:[]}} 
    } = state.httpData;
    return { companiesData, townsData };
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch)
});

export default connect(mapStateToProps, mapDispatchToProps)(UserManager)

