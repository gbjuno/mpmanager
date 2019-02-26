/**
 * Created by Jingle Chen on 2017/12/7.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import * as _ from 'lodash'
import { Table, Button, Row, Col, Calendar, Tabs, message } from 'antd';
import * as CONSTANTS from '../../constants';
import { fetchData, receiveData } from '../../action';
import BreadcrumbCustom from '../../components/BreadcrumbCustom';
import CompanySearch from '../search/CompanySearch';
import * as config from '../../axios/config';
import * as utils from '../../utils';
import moment from 'moment';
import 'moment/locale/zh-cn';
moment.locale('zh-cn');


const TabPane = Tabs.TabPane;
const DIMESION = {
    FULL: 'finish_percentage_all',
    YEAR: 'finish_percentage_last_365_days',
    HALF_A_YEAR: 'finish_percentage_last_182_days',
    QUARTER: 'finish_percentage_last_90_days',
    MONTH: 'finish_percentage_last_30_days',
}


class PhotoStatus extends React.Component {
    state = {
        selectedRowKeys: [],  // Check here to configure the default column
        loading: false,
        companiesData: [],
        selectedCompany: '',
        selectedCompanyId: '',
        selectedRecord:{},
        currentPage: 1,
        visible: false,
        editable: false,
        hasNewRow: false,
        pageSize: 10,
        total: 0,
        expandedRowKeys: [],
    };
    componentDidMount = () => {
        this.start();
    }

    start = () => {
        this.setState({ loading: true });
        this.fetchData();
    };

    fetchData = () => {
        const { fetchData } = this.props
        const { currentPage, pageSize } = this.state
        let tempTownId
        fetchData({
            funcName: 'fetchCompanies', params: {
                
            }, stateName: 'companiesData'
        }).then(res => {
            if (res === undefined || res.data === undefined || res.data.companies === undefined) return
            this.setState({
                companiesData: [...res.data.companies.map(val => {
                    val.key = val.id;
                    return val;
                })],
                total: res.data.count,
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
        if (record.id === -1 || editable) {
            return
        }
        this.setState({
            selectedCompanyId: record.id,
            selectedCompany: record.name,
            selectedRowKeys: selectedRowKeys.length > 0 && selectedRowKeys[0] === record.id ? [] : [record.id],
            selectedRecord: selectedRowKeys.length > 0 && selectedRowKeys[0] === record.id ? {} : record,
        }, () => {
            //this.fetchRelatedUserAndPlace(record.id);
        });

    }


    handlePageChange = (page, pageSize) => {
        this.setState({
            currentPage: page,
        }, () => this.fetchData())
    }

    uploadProps = () => {
        const props =
            {
                name: 'uploadFile',
                action: config.COMPANY_UPLOAD_URL, //TODO: 换成上传地址
                showUploadList: false,
                onChange(info) {
                    if (info.file.status !== 'uploading') {
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

    onTabChange = (key, dimesion) => {
        this.setState({
            [dimesion]: key,
        })
    }

    getColumnsByDimension = (dimesion = DIMESION.FULL) => {
        let columns = [{
            title: '公司名',
            dataIndex: 'name',
            width: "70%",
            render: (text, record) => {
                return <a>{text}</a>
            }
        },
        {
            title: '完成率',
            dataIndex: dimesion,
            width: "30%",
            sorter: (a, b) => {
                return parseFloat(a[[dimesion]]) - parseFloat(b[[dimesion]])
            },
            render: (text, record) => {
                return utils.number(parseFloat(text) * 100, 2) + '%';
            }
        }];

        return columns;
    }

    onPanelChange = (value, mode) => {
        console.log(value, mode);
    }

    render() {
        const { loading, selectedRowKeys, editable, 
            currentPage, pageSize, total, expandedRowKeys,
            selectedRecord } = this.state;
        const { companiesData, countriesC2Data } = this.props
        const rowSelection = {
            selectedRowKeys,
            onChange: this.onSelectChange,
            type: 'radio',
        };

        
        let companiesWrappedData = []
        if (companiesData.data && companiesData.data.companies) {
            companiesWrappedData = [...companiesData.data.companies.map(item => { item.key = item.id; return item })]
        }


        let options = [];
        if (countriesC2Data.data && countriesC2Data.data.countries) {
            options = [...countriesC2Data.data.countries.map(item => { item.key = item.id; return item })]
        }

        const hasSelected = selectedRowKeys.length > 0 && selectedRowKeys[0] !== -1



        return (
            <div className="gutter-example">
                <BreadcrumbCustom first="完成率统计" />
                <CompanySearch fetchData={fetchData} />
                <Row gutter={16}>
                <Col className="gutter-row" md={14}>
                    <Calendar onPanelChange={this.onPanelChange} />
                </Col>
                </Row>
            </div>
        );
    }
}

const mapStateToProps = state => {
    const {
        companiesData = { data: { count: 0, towns: [] } },
        townsData = { data: { count: 0, towns: [] } },
        countriesC2Data = { data: { count: 0, countries: [] } },
        usersInCompany = { data: { count: 0, users: [] } },
        placesInCompany = { data: { count: 0, monitor_places: [] } },
    } = state.httpData;
    return { companiesData, townsData, countriesC2Data, usersInCompany, placesInCompany };
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch)
});

export default connect(mapStateToProps, mapDispatchToProps)(PhotoStatus)

