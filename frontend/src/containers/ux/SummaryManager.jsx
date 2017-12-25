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
        this.fetchData();
    };

    fetchData = () => {
        const { fetchData } = this.props
        const { currentPage } = this.state
        let tempTownId
        fetchData({funcName: 'fetchSummaries',params: {pageNo: currentPage, pageSize: 20}, stateName: 'summariesData'}).then(res => {
            if(res === undefined || res.data === undefined || res.data.summaries === undefined) return
            this.setState({
                summariesData: [...res.data.summaries.map(val => {
                    val.key = val.company_id;
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

    render() {
        const summaryColumns = [{
                title: '统计日期',
                dataIndex: 'day',
                width: 20,
                render: (text, record) => {
                    if (record.id === -1){
                        return ''
                    }
                    var createAt = new Date(text).toLocaleString('chinese',{hour12:false});
                    return createAt.substring(0, createAt.indexOf(' '))
                }
            },{
            title: '公司',
            dataIndex: 'company_name',
            width: 20,
            render: (text, record) => {
                return <a href={record.url} target="_blank">{text}</a>                
            }
        }, {
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
        }, {
            title: '未完成的拍照地点',
            dataIndex: 'unfinish_ids',
            width: 40,
            render: (text, record) => {
                return <a href={record.url} target="_blank">{text}</a>
            }
        }];

        const { loading, selectedRowKeys, selectedTown, summariesData } = this.state;
        console.log('summaries ...sss...', summariesData)
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
                <BreadcrumbCustom first="安监管理" second="统计报表" />
                <Row gutter={16}>
                    <Col className="gutter-row" md={24}>
                        <div className="gutter-box">
                            <Card title="统计报告" bordered={false}>
                                <Table columns={summaryColumns} dataSource={summariesData}
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

export default connect(mapStateToProps, mapDispatchToProps)(SummaryManager)

