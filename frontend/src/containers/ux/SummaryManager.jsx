/**
 * Created by Jingle Chen on 2017/12/7.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { Table, Button, Row, Col, Card, Input, Icon } from 'antd';
import moment from 'moment';
import * as CONSTANTS from '../../constants';
import { fetchData, receiveData } from '../../action';
import { getPros } from '../../axios';
import BreadcrumbCustom from '../../components/BreadcrumbCustom';
import SummarySearch from '../search/SummarySearch';



class SummaryManager extends React.Component {
    state = {
        selectedRowKeys: [],  // Check here to configure the default column
        loading: false,
        summariesData: [],
        selectedSummary: '',
        selectedSummaryId: '',
        currentPage: 1
    };

    componentDidMount() {
        this.start();
    }

    start = () => {
        this.setState({ loading: true });
        //this.fetchData();
        this.searchSummary(moment(new Date()).format(CONSTANTS.DATE_QUERY_FORMAT))
    };

    fetchData = () => {
        const { fetchData } = this.props
        const { currentPage } = this.state
        let tempTownId
        fetchData({funcName: 'fetchSummaries',params: {pageNo: currentPage, pageSize: 20}, stateName: 'summariesData'}).then(res => {
            if(res === undefined || res.data === undefined || res.data.summaries === undefined) return
            this.setState({
                summariesData: [...res.data.summaries.map(val => {
                    val.key = val.id;
                    return val;
                })],
                loading: false,
            });
        });
    }

    searchSummary = (date) => {
        const { fetchData } = this.props
        fetchData({funcName: 'searchSummaries', params: {day:date}, 
            stateName: 'summariesData'})
    }

    onSelectChange = (selectedRowKeys) => {
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

    render() {
        const summaryColumns = [
            {
                title: '统计日期',
                dataIndex: 'day',
                width: 20,
                render: (text, record) => {
                    if (record.id === -1){
                        return ''
                    }
                    return moment(new Date(text)).format(CONSTANTS.DATE_DISPLAY_FORMAT)
                }
            },{
                title: '公司',
                dataIndex: 'company_name',
                width: 20,
                render: (text, record) => {
                    return <a href={record.url} target="_blank">{text}</a>
                }
            },{
                title: '是否完成',
                dataIndex: 'finish',
                width: 20,
                render: (text, record) => {
                    var value 
                    if(text == "T") {
                        value = "是"
                    } else {
                        value = "否"
                    }
                    return <a href={record.url} target="_blank">{value}</a>
                }
            },{
                title: '未完成的拍照地点',
                dataIndex: 'unfinish_ids',
                width: 40,
                render: (text, record) => {
                    return <a href={record.url} target="_blank">{text}</a>
                }
            }
        ];

        const { loading, selectedRowKeys, selectedTown } = this.state;
        const { summariesData } = this.props
        const rowSelection = {
            selectedRowKeys,
            onChange: this.onSelectChange,
        };

        let summariesWrappedData = []
        if(summariesData.data && summariesData.data.summaries){
            summariesWrappedData = [...summariesData.data.summaries.map(item => {item.key = item.id; return item})]
        }

        return (
            <div className="gutter-example">
                <BreadcrumbCustom first="统计报表" second="" />
                <SummarySearch  fetchData={fetchData}/>
                <Row gutter={16}>
                    <Col className="gutter-row" md={24}>
                        <div className="gutter-box">
                            <Card title="统计报告" bordered={false}>
                                <Table columns={summaryColumns} dataSource={summariesWrappedData}
                                    size="small" onRowClick={this.onRowClick}
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
        summariesData = {data: {count:0, towns:[]}}, 
    } = state.httpData;
    return { summariesData };
};

const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch)
});

export default connect(mapStateToProps, mapDispatchToProps)(SummaryManager)

