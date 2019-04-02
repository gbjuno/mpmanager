/**
 * Created by Jingle Chen on 2017/12/7.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import moment from 'moment';
import * as _ from 'lodash'
import * as download from 'downloadjs'
import { Table, Button, Row, Col, Card, Input, Icon, Pagination, Modal, Upload, Tabs, message } from 'antd';
import * as CONSTANTS from '../../constants';
import { fetchData, receiveData } from '../../action';
import BreadcrumbCustom from '../../components/BreadcrumbCustom';
import EditableCell from '../../components/cells/EditableCell';
import PhotoStatusSearch from '../search/PhotoStatusSearch';
import PhotoStatusDownloader from '../search/PhotoStatusDownloader';
import * as config from '../../axios/config';
import { Bar } from '../../components/Charts';
import DataSet from "@antv/data-set";
import * as utils from '../../utils';
import styles from '../less/PhotoStatus.less';

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


    downloadFile = (year, month) => {

        let url = config.COMPANY_REPORT_EXPORT_URL(year, month); 
        const x = new XMLHttpRequest;
        x.open("GET", url, true);
        x.responseType = "blob";
        x.withCredentials = true;
        x.onload = function (e) {
            download(x.response, "完成情况报表.xlsx", "application/octet-stream")
        }
        x.send();
    }



    // The upper is place operations

    getPhotoStatus = (selectedRecord) => {
        if(_.isEmpty(selectedRecord)) return {}

        const { DataView } = DataSet;
        const dv = new DataView();

        let selectedComData = selectedRecord
        let photoStatusData = [
            {
                x: '全部完成率',
                y: selectedComData.finish_percentage_all,
            },
            {
                x: '年完成率',
                y: selectedComData.finish_percentage_last_365_days,
            },
            {
                x: '半年完成率',
                y: selectedComData.finish_percentage_last_182_days,
            },
            {
                x: '季完成率',
                y: selectedComData.finish_percentage_last_90_days,
            },
            {
                x: '月完成率',
                y: selectedComData.finish_percentage_last_30_days,
            },
        ]
        dv.source(photoStatusData);

        return dv
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

        let photoStatusData = this.getPhotoStatus(selectedRecord)


        let options = [];
        if (countriesC2Data.data && countriesC2Data.data.countries) {
            options = [...countriesC2Data.data.countries.map(item => { item.key = item.id; return item })]
        }

        const hasSelected = selectedRowKeys.length > 0 && selectedRowKeys[0] !== -1

        const dimesionList = [{
            key: DIMESION.FULL,
            tab: '全部',
          }, {
            key: DIMESION.YEAR,
            tab: '最近一年',
          }, {
            key: DIMESION.HALF_A_YEAR,
            tab: '最近半年',
          }, {
            key: DIMESION.QUARTER,
            tab: '最近一季',
          }, {
            key: DIMESION.MONTH,
            tab: '最近一月',
          }];

        let companyColumns = this.getColumnsByDimension(this.state.noTitleKey)


        return (
            <div className="gutter-example">
                <BreadcrumbCustom first="完成率统计" />
                <PhotoStatusSearch fetchData={fetchData} />
                <PhotoStatusDownloader fetchData={fetchData} exportReport={this.downloadFile}/>
                <Row gutter={16}>
                <Col className="gutter-row" md={14}>
                    <Card bordered={false}
                         tabList={dimesionList}
                         onTabChange={(key) => { this.onTabChange(key, 'noTitleKey'); }}
                    >
                        <div className="gutter-box">
                                <Table rowSelection={rowSelection}
                                    size="small"
                                    columns={companyColumns}
                                    dataSource={companiesWrappedData}
                                    onRow={(record) => ({
                                        onClick: () => this.onRowClick(record),
                                    })}
                                    pagination={{
                                        hideOnSinglePage: true,
                                        onChange: this.handlePageChange,
                                        current: currentPage,
                                        defaultCurrent: 1,
                                        pageSize,
                                        // total,
                                    }}
                                />
                        </div>
                        </Card>
                    </Col>
                    <Col className="gutter-row" md={9}>
                        <Card bordered={false}>
                        <div className="gutter-box" style={{height: 400}}>
                        {hasSelected?
                        <Bar
                            height={400}
                            title={
                            "完成率"
                            }
                            data={photoStatusData}
                        />
                        :
                        <p>点击公司查看详情</p>
                        }
                        </div>
                        </Card>
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

