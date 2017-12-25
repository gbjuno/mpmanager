/**
 * Created by Jingle Chen on 2017/12/7.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { Table, Button, Row, Col, Card, Input, Icon } from 'antd';
import { fetchData, receiveData } from '../../action';
import { getPros } from '../../axios';
import BreadcrumbCustom from '../../components/BreadcrumbCustom';



class EditableCell extends React.Component {

    state = {
        value: this.props.value,
        editable: true,
    }
    handleChange = (e) => {
        const value = e.target.value;
        this.setState({ value });
    }
    check = () => {
        this.setState({ editable: false });
        if (this.props.onChange) {
            this.props.onChange(this.state.value);
        }
    }
    edit = () => {
        this.setState({ editable: true });
    }

    render(){
        const { value, editable } = this.state;
        return (
            <div className="editable-cell">
                {
                editable ?
                    <div className="editable-cell-input-wrapper">
                    <Input
                        value={value}
                        onChange={this.handleChange}
                        onPressEnter={this.check}
                    />
                    {/** 
                    <Icon
                        type="check"
                        className="editable-cell-icon-check"
                        onClick={this.check}
                    />
                    */}
                    </div>
                    :
                    <a target="_blank">{value}</a>
                }
            </div>
        )
    }
}

class UserManager extends React.Component {
    state = {
        selectedRowKeys: [],  // Check here to configure the default column
        loading: false,
        usersData: [],
        selectedCompany: '',
        selectedCompanyId: '',
    };
    componentDidMount() {
        this.start();
    }
    start = () => {
        this.setState({ loading: true });
        this.fetchData();
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


    onSelectChange = (selectedRowKeys) => {
        console.log('selectedRowKeys changed: ', selectedRowKeys);
        if(selectedRowKeys.length > 0){
            selectedRowKeys = [selectedRowKeys[selectedRowKeys.length-1]]
        }
        
        this.setState({ selectedRowKeys });
    };

    onRowClick = (record, index, event) => {
        const { selectedRowKeys } = this.state
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
                address: '',
                country_id: '',
            }, ...this.state.usersData]
        });
    }

    handleDeleteTown = () => {
        const { fetchData } = this.props
        const { townSelectedRowKeys } = this.state
        if(townSelectedRowKeys.length === 0) return
        fetchData({funcName: 'deleteTown', params: {townId: townSelectedRowKeys[0]}, stateName: 'deleteTownStatus'})
            .then(res => {
                console.log('delete town successfully', res)
                this.fetchTownsData() 
            }).catch(err => console.log(err));
    }

    onNewTownSave = (name) => {
        const { fetchData } = this.props
        fetchData({funcName: 'newTown', params: {name}, stateName: 'newTownStatus'})
            .then(res => {
                console.log('create new town successfully', res)
                this.fetchTownsData() 
            }).catch(err => console.log(err));
        console.log('value--->', name)
    }

    render() {

        const userColumns = [{
            title: '用户名',
            dataIndex: 'name',
            width: 40,
            render: (text, record) => {
                if(record.id === -1){
                    return <EditableCell value={record.name} onChange={this.onNewTownSave} />
                }else{
                    return <a href={record.url} target="_blank">{text}</a>
                }
            }
        }, {
            title: '手机号',
            dataIndex: 'phone',
            width: 40,
            render: (text, record) => {
                if(record.id === -1){
                    return <EditableCell value={record.country_id} onChange={this.onNewTownSave} />
                }else{
                    return <a href={record.url} target="_blank">{text}</a>
                }
            }
        }, {
            title: '职位',
            dataIndex: 'job',
            width: 40,
            render: (text, record) => {
                if(record.id === -1){
                    return <EditableCell value={record.address} onChange={this.onNewTownSave} />
                }else{
                    return <a href={record.url} target="_blank">{text}</a>
                }
            }
        }, {
            title: '所在公司',
            dataIndex: 'company_name',
            width: 40,
            render: (text, record) => {
                if(record.id === -1){
                    return <EditableCell value={record.address} onChange={this.onNewTownSave} />
                }else{
                    return <a href={record.url} target="_blank">{text}</a>
                }
            }
        },{
            title: '微信号',
            dataIndex: 'wx_openid',
            width: 80,
            render: (text, record) => {
                return <a href={record.url} target="_blank">{text}</a>
            }
        }];

        const { loading, selectedRowKeys, selectedTown,
            usersData } = this.state;
        const rowSelection = {
            selectedRowKeys,
            onChange: this.onSelectChange,
        };
        return (
            <div className="gutter-example">
                <style>
                {`
                    .editable-cell {
                        position: relative;
                      }
                      
                      .editable-cell-input-wrapper,
                      .editable-cell-text-wrapper {
                        padding-right: 24px;
                      }
                      
                      .editable-cell-text-wrapper {
                        padding: 5px 24px 5px 5px;
                      }
                      
                      .editable-cell-icon,
                      .editable-cell-icon-check {
                        position: absolute;
                        right: 0;
                        width: 20px;
                        cursor: pointer;
                      }
                      
                      .editable-cell-icon {
                        line-height: 18px;
                        display: none;
                      }
                      
                      .editable-cell-icon-check {
                        line-height: 28px;
                      }
                      
                      .editable-cell:hover .editable-cell-icon {
                        display: inline-block;
                      }
                      
                      .editable-cell-icon:hover,
                      .editable-cell-icon-check:hover {
                        color: #108ee9;
                      }
                      
                      .editable-add-btn {
                        margin-bottom: 8px;
                      }
                `}
                </style>
                <BreadcrumbCustom first="安监管理" second="用户管理" />
                <Row gutter={16}>
                    <Col className="gutter-row" md={24}>
                        <div className="gutter-box">
                            <Card title="用户列表" bordered={false}>
                                <div style={{ marginBottom: 16 }}>
                                    <Button type="primary" onClick={this.handleAdd}
                                            disabled={loading} 
                                    >新增</Button>
                                    <Button type="primary" onClick={this.handleDeleteTown}
                                            disabled={loading} 
                                    >删除</Button>
                                </div>
                                <Table rowSelection={rowSelection} columns={userColumns} dataSource={usersData}
                                        onRowClick={this.onRowClick}
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
        townsData = {data: {count:0, towns:[]}}, 
        fetchCountries = {data: {count:0, countries:[]}} 
    } = state.httpData;
    return { townsData };
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch)
});

export default connect(mapStateToProps, mapDispatchToProps)(UserManager)

