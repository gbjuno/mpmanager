/**
 * Created by Jingle Chen on 2017/12/7.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import * as _ from 'lodash'
import { Table, Button, Row, Col, Card, Input, Icon, message } from 'antd';
import { fetchData, receiveData, searchFilter, resetFilter } from '../../action';
import { getPros } from '../../axios';
import BreadcrumbCustom from '../../components/BreadcrumbCustom';
import EditableCell from '../../components/cells/EditableCell';
import UserSearch from '../search/UserSearch';



class UserManager extends React.Component {
    state = {
        selectedRowKeys: [],  // Check here to configure the default column
        loading: false,
        usersData: [],
        companiesData: [],
        selectedCompany: '',
        selectedCompanyId: '',
        editable: false,
        hasNewRow: false,
        currentPage: 1,
        pageSize: 10,
        total: 0,
    };
    componentDidMount() {
        this.start();
    }
    start = () => {
        const { resetFilter, searchFilter } = this.props
        const { currentPage, pageSize } = this.state
        resetFilter('user')
        searchFilter('user', {
            pageNo: currentPage,
            pageSize: pageSize,
        })
        this.setState({ loading: true }, () => {
            this.fetchData();
            this.fetchCompanyList();
        });
        
    };

    fetchData = () => {
        const { fetchData, searchFilter, filter } = this.props

        fetchData({funcName: 'fetchUsers', params: filter.user, stateName: 'usersData'}).then(res => {
            if(res === undefined || res.data === undefined || res.data.users === undefined) return
            this.setState({
                usersData: [...res.data.users.map(val => {
                    val.key = val.id;
                    return val;
                })],
                loading: false,
            });
            searchFilter('user', {
                total: res.data.count,
            })
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
            hasNewRow: true,
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
            hasNewRow: false,
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

    handlePageChange = (page, pageSize) => {
        const { searchFilter } = this.props
        searchFilter('user', {
            pageSize: 10,
            pageNo: page,
        })
        this.setState({
            currentPage: page,
        }, () => this.fetchData())
    }

    render() {

        const { loading, selectedRowKeys, companiesData, hasNewRow, pageSize,
             editable } = this.state;
        const { usersData, filter } = this.props;
        const rowSelection = {
            selectedRowKeys,
            onChange: this.onSelectChange,
            type: 'radio',
        };

        let total = 0; 
        let currentPage = 1;
        if(filter.user) {
            total = filter.user.total
            currentPage = filter.user.pageNo
        }
        console.log('sally, i will take all my life to protect you...', filter)

        let options = [];
        if(companiesData){
            options = [...companiesData.map(item => {item.key = item.id; return item})]
        }

        let usersWrappedData = []
        if(usersData.data && usersData.data.users){
            usersWrappedData = [...usersData.data.users.map(item => {item.key = item.id; return item})]
        }

        if(hasNewRow){
            usersWrappedData = [{
                key: -1,
                id: -1,
                name: '',
                phone: '',
                job: '',
                company_id: '',
            }, ...usersWrappedData]
        }else{
            usersWrappedData = [...usersWrappedData.filter(item => item.id !== -1)]
        }

        const hasSelected = selectedRowKeys.length > 0 && selectedRowKeys[0] !== -1

        const userColumns = [{
            title: '用户名',
            dataIndex: 'name',
            width: '15%',
            render: (text, record) => {
                if (record.id === -1 || (editable && record.id === selectedRowKeys[0])) {
                    return <EditableCell dataIndex="user.name" value={record.name} onChange={this.onNewRowChange} />
                }
                return <a>{text}</a>
            }
        }, {
            title: '手机号',
            dataIndex: 'phone',
            width: '15%',
            render: (text, record) => {
                if (record.id === -1 || (editable && record.id === selectedRowKeys[0])) {
                    return <EditableCell dataIndex="user.phone" value={record.phone} onChange={this.onNewRowChange} />
                }
                return <a>{text}</a>
            }
        }, {
            title: '职位',
            dataIndex: 'job',
            width: '15%',
            render: (text, record) => {
                if (record.id === -1 || (editable && record.id === selectedRowKeys[0])) {
                    return <EditableCell dataIndex="user.job" value={record.job} onChange={this.onNewRowChange} />
                }
                return <a>{text}</a>
            }
        }, {
            title: '所在公司',
            dataIndex: 'company_name',
            width: '25%',
            render: (text, record) => {
                if (record.id === -1 || (editable && record.id === selectedRowKeys[0])) {
                    return <EditableCell dataIndex="user.company_id" value={record.company_id} onChange={this.onNewRowChange}
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
                <BreadcrumbCustom first="用户管理" second="" />
                <UserSearch fetchData={fetchData} />
                <Row gutter={16}>
                    <Col className="gutter-row" md={24}>
                        <div className="gutter-box">
                            <Card title="用户列表" bordered={false}>
                                <div style={{ marginBottom: 16 }}>
                                    <Button type="primary" onClick={this.handleAdd}
                                    >新增</Button>
                                    <Button type="primary" onClick={this.handleModify}
                                        disabled={!hasSelected}
                                    >修改</Button>
                                    <Button type="primary" onClick={this.handleDelete}
                                        disabled={!hasSelected}
                                    >删除</Button>
                                </div>
                                <Table rowSelection={rowSelection} columns={userColumns} dataSource={usersWrappedData}
                                    size="small"
                                    onRow={(record) => ({
                                        onClick: () => this.onRowClick(record),
                                    })}
                                    pagination={{
                                        hideOnSinglePage: true,
                                        onChange: this.handlePageChange,
                                        current: currentPage,
                                        defaultCurrent: 1,
                                        pageSize,
                                        total,
                                    }}
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
        usersData = {data: {count:0, companies:[]}}, 
    } = state.httpData;
    return { usersData, filter: state.searchFilter };
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch),
    searchFilter: bindActionCreators(searchFilter, dispatch),
    resetFilter: bindActionCreators(resetFilter, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(UserManager)

