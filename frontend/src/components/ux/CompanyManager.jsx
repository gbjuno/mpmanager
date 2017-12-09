/**
 * Created by Jingle Chen on 2017/12/7.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { Table, Button, Row, Col, Card, Input, Icon, Pagination } from 'antd';
import { fetchData, receiveData } from '../../action';
import { getPros } from '../../axios';
import BreadcrumbCustom from '../BreadcrumbCustom';



class EditableCell extends React.Component {

    state = {
        dataIndex : this.props.dataIndex,
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

class CompanyManager extends React.Component {
    state = {
        selectedRowKeys: [],  // Check here to configure the default column
        loading: false,
        companiesData: [],
        selectedCompany: '',
        selectedCompanyId: '',
        currentPage: 1
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
        const { currentPage } = this.state
        let tempTownId
        fetchData({funcName: 'fetchCompanies', params: {
                pageNo: currentPage, pageSize: 20}, stateName: 'companiesData'}).then(res => {
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
            companiesData: [{
                key: -1,
                id: -1,
                name: '',
                address: '',
                country_id: '',
            }, ...this.state.companiesData]
        });
    }

    handleDeleteTown = () => {
        const { fetchData } = this.props
        const { townSelectedRowKeys, currentPage } = this.state
        if(townSelectedRowKeys.length === 0) return
        fetchData({funcName: 'deleteTown', params: {townId: townSelectedRowKeys[0], 
            pageNo: currentPage, pageSize: 20}, stateName: 'deleteTownStatus'})
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

    getPagination = () => {
        return <Pagination onChange={this.handlePageChange}/>
    }

    handlePageChange = (page, pageSize) => {
        this.setState({
            currentPage: page,
        })
    }

    render() {

        const companyColumns = [{
            title: '公司名',
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
            title: '所在镇',
            dataIndex: 'country_id',
            width: 40,
            render: (text, record) => {
                if(record.id === -1){
                    return <EditableCell value={record.country_id} onChange={this.onNewTownSave} />
                }else{
                    return <a href={record.url} target="_blank">{text}</a>
                }
            }
        }, {
            title: '详细地址',
            dataIndex: 'address',
            width: 40,
            render: (text, record) => {
                if(record.id === -1){
                    return <EditableCell value={record.address} onChange={this.onNewTownSave} />
                }else{
                    return <a href={record.url} target="_blank">{text}</a>
                }
            }
        }, {
            title: '创建时间',
            dataIndex: 'create_at',
            width: 80,
            render: (text, record) => {
                if (record.id === -1){
                    return ''
                }
                var createAt = new Date(text).toLocaleString('chinese',{hour12:false});
                return createAt;
            }
        }];

        const { loading, selectedRowKeys, selectedTown,
            companiesData } = this.state;
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
                <BreadcrumbCustom first="安监管理" second="公司管理" />
                <Row gutter={16}>
                    <Col className="gutter-row" md={24}>
                        <div className="gutter-box">
                            <Card title="公司列表" bordered={false}>
                                <div style={{ marginBottom: 16 }}>
                                    <Button type="primary" onClick={this.handleAdd}
                                            disabled={loading} 
                                    >新增</Button>
                                    <Button type="primary" onClick={this.handleDeleteTown}
                                            disabled={loading} 
                                    >删除</Button>
                                </div>
                                <Table rowSelection={rowSelection} columns={companyColumns} dataSource={companiesData}
                                        onRowClick={this.onRowClick} pagination={this.getPagination()}
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

export default connect(mapStateToProps, mapDispatchToProps)(CompanyManager)

