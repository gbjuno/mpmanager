/**
 * Created by Jingle Chen on 2017/12/7.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import moment from 'moment';
import * as _ from 'lodash'
import * as download from 'downloadjs'
import { Table, Button, Row, Col, Card, Input, Icon, Pagination, Modal, Upload, message } from 'antd';
import * as CONSTANTS from '../../constants';
import { fetchData, receiveData } from '../../action';
import { getPros } from '../../axios';
import BreadcrumbCustom from '../../components/BreadcrumbCustom';
import EditableCell from '../../components/cells/EditableCell';
import CompanySearch from './search/CompanySearch';


class CompanyManager extends React.Component {
    state = {
        selectedRowKeys: [],  // Check here to configure the default column
        loading: false,
        companiesData: [],
        selectedCompany: '',
        selectedCompanyId: '',
        currentPage: 1,
        visible: false,
    };
    componentDidMount() {
        this.start();
    }
    start = () => {
        this.setState({ loading: true });
        this.fetchData();
        this.fetchCountryListWithoutTownId();
    };

    fetchData = () => {
        const { fetchData } = this.props
        const { currentPage } = this.state
        let tempTownId
        fetchData({
            funcName: 'fetchCompanies', params: {
                pageNo: currentPage, pageSize: 20
            }, stateName: 'companiesData'
        }).then(res => {
            if (res === undefined || res.data === undefined || res.data.companies === undefined) return
            this.setState({
                companiesData: [...res.data.companies.map(val => {
                    val.key = val.id;
                    return val;
                })],
                loading: false,
            });
        });
    }

    fetchCountryListWithoutTownId = () => {
        const { fetchData } = this.props
        fetchData({funcName: 'fetchCountriesWithoutTownId', stateName: 'countriesC2Data', params: {}}).then(res => {
            if(res === undefined || res.data === undefined || res.data.countries === undefined) return
            this.setState({
                countriesData: [...res.data.countries.map(val => {
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
        console.log('selectedRowKeys changed: ', selectedRowKeys);
        if (selectedRowKeys.length > 0) {
            selectedRowKeys = [selectedRowKeys[selectedRowKeys.length - 1]]
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

    handleCancelEditRow = () => {
        let tmpCompaniesDataa = [...this.state.companiesData.filter(item => item.id !== -1)]
        this.setState({
            companiesData: tmpCompaniesDataa
        })
    }

    handleDelete = () => {
        const { fetchData } = this.props
        const { selectedRowKeys, currentPage } = this.state
        if (selectedRowKeys.length === 0) return
        fetchData({
            funcName: 'deleteCompany', params: { id: selectedRowKeys[0] }, stateName: 'deleteCompanyStatus'
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

    onNewRowSave = () => {
        const { fetchData } = this.props
        const keys = _.keys(this.state)
        const PREFIX = 'company.'
        const PREFIX_LEN = PREFIX.length;
        let obj = {}
        for (let key of keys) {
            if(_.startsWith(key, PREFIX)){
                let field = key.substring(PREFIX_LEN)
                obj[field] = this.state[key]
            }
        }

        fetchData({funcName: 'newCompany', params: obj, stateName: 'newCompanyStatus'})
            .then(res => {
                message.success('创建成功')
                this.fetchData()
            }).catch(err => {
                let errRes = err.response
                if(errRes.data && errRes.data.status === 'error'){
                    message.error(errRes.data.error)
                }
            });
    }

    getPagination = () => {
        return <Pagination onChange={this.handlePageChange} />
    }

    handlePageChange = (page, pageSize) => {
        this.setState({
            currentPage: page,
        })
    }

    uploadProps = () => {
        const props = 
        {
            name: 'file',
            action: '//jsonplaceholder.typicode.com/posts/', //TODO: 换成上传地址
            headers: {
            authorization: 'authorization-text',
            },
            onChange(info) {
                if (info.file.status !== 'uploading') {
                    console.log(info.file, info.fileList);
                }
                if (info.file.status === 'done') {
                    message.success(`${info.file.name}上传成功`);
                } else if (info.file.status === 'error') {
                    message.error(`${info.file.name}上传失败`);
                }
            },
        }
        return props
    }

    downloadFile = () => {
        let url = 'http://localhost:3006/static/media/b1.25566666.png'; //TODO: 换成下载公司数据url,及相应的文件格式
        download(`${url}`)
    }

    /**
     * 准备删
     */
    showModal = () => {
        this.setState({
            visible: true,
        })
    }

    hideModal = () => {
        this.setState({
            visible: false,
        })
    }

    render() {

        const { loading, selectedRowKeys, selectedTown,
            companiesData } = this.state;
        const { countriesC2Data } = this.props
        const rowSelection = {
            selectedRowKeys,
            onChange: this.onSelectChange,
            type: 'radio',
        };

        let options = [];//[...countriesC2Data.map(item => item.key = item.id)];
        if(countriesC2Data.data && countriesC2Data.data.countries){
            options = [...countriesC2Data.data.countries.map(item => {item.key = item.id; return item})]
        }

        const hasSelected = selectedRowKeys.length > 0

        const companyColumns = [{
            title: '公司名',
            dataIndex: 'name',
            width: "20%",
            render: (text, record) => {
                if (record.id === -1) {
                    return <EditableCell dataIndex='company.name' value={record.name} onChange={this.onNewRowChange} />
                }
                return <a>{text}</a>
            }
        }, {
            title: '所在村',
            dataIndex: 'country_name',
            width: "20%",
            render: (text, record) => {
                if (record.id === -1) {
                    return <EditableCell dataIndex='company.country_id' value={record.country_id} onChange={this.onNewRowChange}
                    editType="select" valueType="int" options={options} placeholder="请选择村"/>
                }
                return <a>{text}</a>
            }
        }, {
            title: '详细地址',
            dataIndex: 'address',
            width: "30%",
            render: (text, record) => {
                if (record.id === -1) {
                    return <EditableCell dataIndex='company.address' value={record.address} onChange={this.onNewRowChange} />
                }
                return <a href={record.url} target="_blank">{text}</a>
            }
        }, {
            title: '创建时间',
            dataIndex: 'create_at',
            width: "30%",
            render: (text, record) => {
                if (record.id === -1) {
                    return <EditableCell type="opt" onSave={this.onNewRowSave} onCancel={this.handleCancelEditRow}/>
                }
                var createAt = moment(new Date(text)).format(CONSTANTS.DATE_TABLE_FORMAT)
                return createAt;
            }
        }];

        
        return (
            <div className="gutter-example">
                <BreadcrumbCustom first="安监管理" second="公司管理" />
                <CompanySearch  fetchData={fetchData}/>
                <Row gutter={16}>
                    <Col className="gutter-row" md={24}>
                        <div className="gutter-box">
                            <Card title="公司列表" bordered={false}>
                                <div style={{ marginBottom: 16 }}>
                                    <Button type="primary" onClick={this.handleAdd}
                                        disabled={loading}
                                    >新增</Button>
                                    <Button type="primary" onClick={this.showModal}
                                        disabled={!hasSelected}
                                    >修改</Button>
                                    <Modal
                                        title="还没做"
                                        visible={this.state.visible}
                                        onOk={this.hideModal}
                                        onCancel={this.hideModal}
                                        okText="确认"
                                        cancelText="取消"
                                        >
                                        <p>别点了，没有做</p>
                                    </Modal>
                                    <Button type="primary" onClick={this.handleDelete}
                                        disabled={!hasSelected}
                                    >删除</Button>
                                    <Button type="primary" onClick={this.downloadFile}
                                        disabled={loading}
                                    >下载</Button>
                                    <Upload style={{marginLeft: '10px'}} {...this.uploadProps}>
                                        <Button type="primary">上传
                                        </Button>
                                    </Upload>
                                    
                                </div>
                                <Table rowSelection={rowSelection} columns={companyColumns} dataSource={companiesData}
                                    onRow={(record) => ({
                                        onClick: () => this.onRowClick(record),
                                    })}
                                    pagination={this.getPagination()}
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
        townsData = { data: { count: 0, towns: [] } },
        countriesData = { data: { count: 0, countries: [] } },
        countriesC2Data = { data: { count: 0, countries: [] } },
    } = state.httpData;
    return { townsData, countriesC2Data };
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch)
});

export default connect(mapStateToProps, mapDispatchToProps)(CompanyManager)

