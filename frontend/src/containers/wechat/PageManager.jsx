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
import { getPros } from '../../axios';
import BreadcrumbCustom from '../../components/BreadcrumbCustom';
import EditableCell from '../../components/cells/EditableCell';
import * as config from '../../axios/config';

const TabPane = Tabs.TabPane;

class PageManager extends React.Component {
    state = {
        selectedRowKeys: [],  // Check here to configure the default column
        loading: false,
        pagesData: [],
        selectedPage: '',
        selectedPageId: '',
        currentPage: 1,
        visible: false,
        editable: false,
        hasNewRow: false,
        pageSize: 10,
        total: 0,
        expandedRowKeys: [],

        chapterEditable: false,
        selectedChapterKeys: [],
        selectedChapterId: '',
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
            funcName: 'fetchPages', params: {
                pageNo: currentPage, pageSize: pageSize
            }, stateName: 'pagesData'
        }).then(res => {
            if (res === undefined || res.data === undefined || res.data.templatePages === undefined) return
            this.setState({
                pagesData: [...res.data.templatePages.map(val => {
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
            selectedPageId: record.id,
            selectedPage: record.name,
            selectedRowKeys: selectedRowKeys.length > 0 && selectedRowKeys[0] === record.id ? [] : [record.id],
        });
    }

    handleAdd = () => {
        this.setState({
            currentPage: 1,
            hasNewRow: true,
        });
    }

    handleCancelEditRow = () => {
        let tmpPagesData = [...this.state.pagesData.filter(item => item.id !== -1)]
        this.setState({
            editable: false,
            hasNewRow: false,
            pagesData: tmpPagesData,
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
            funcName: 'deletePage', params: { id: selectedRowKeys[0] }, stateName: 'deletePageStatus'
        }).then(res => {
            message.success('删除成功')
            this.setState({
                visible: false,
                selectedPage: '',
                selectedPageId: '',
            })
            this.fetchData()
        }).catch(err => {
            let errRes = err.response
            if (errRes.data && errRes.data.status === 'error') {
                message.error(errRes.data.error)
                this.setState({
                    visible: false,
                })
            }
        });
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

    onRowSave = () => {
        const { editable } = this.state
        if (editable) {
            this.onUpdateRowSave()
        } else {
            this.onNewRowSave()
        }
    }

    onUpdateRowSave = () => {
        const { fetchData } = this.props
        const { selectedRowKeys } = this.state
        const keys = _.keys(this.state)
        const PREFIX = 'templatepage.'
        const PREFIX_LEN = PREFIX.length;
        let obj = {}
        for (let key of keys) {
            if (_.startsWith(key, PREFIX)) {
                let field = key.substring(PREFIX_LEN)
                obj[field] = this.state[key]
            }
        }
        obj.id = selectedRowKeys[0]

        fetchData({ funcName: 'updatePage', params: obj, stateName: 'updatePageStatus' })
            .then(res => {
                message.success('更新成功')
                this.fetchData()
                this.setState({
                    editable: false,
                })
            }).catch(err => {
                let errRes = err.response
                if (errRes.data && errRes.data.status === 'error') {
                    message.error(errRes.data.error)
                }
            });
    }

    onNewRowSave = () => {
        const { fetchData } = this.props
        const keys = _.keys(this.state)
        const PREFIX = 'templatepage.'
        const PREFIX_LEN = PREFIX.length;
        let obj = {}
        for (let key of keys) {
            if (_.startsWith(key, PREFIX)) {
                let field = key.substring(PREFIX_LEN)
                obj[field] = this.state[key]
            }
        }

        fetchData({ funcName: 'newPage', params: obj, stateName: 'newPageStatus' })
            .then(res => {
                message.success('创建成功')
                this.fetchData()
                this.setState({
                    editable: false,
                    hasNewRow: false,
                })
            }).catch(err => {
                let errRes = err.response
                if (errRes.data && errRes.data.status === 'error') {
                    message.error(errRes.data.error)
                }
            });
    }

    handlePageChange = (page, pageSize) => {
        this.setState({
            currentPage: page,
        }, () => this.fetchData())
    }


    onChapterChange = (selectedChapterKeys) => {
        if (selectedChapterKeys.length > 0) {
            selectedChapterKeys = [selectedChapterKeys[selectedChapterKeys.length - 1]]
        }

        this.setState({ selectedChapterKeys });
    };

    onChapterSelect = (record, selected, selectedChapterKeys) => {
        if (selected && selectedChapterKeys[0] === record.id) {
            this.setState({
                selectedChapterKeys: []
            })
        }
    }

    onChapterClick = (record, index, event) => {
        const { selectedChapterKeys, chapterEditable } = this.state
        if (record.id === -1 || chapterEditable) {
            return
        }
        this.setState({
            selectedChapterId: record.id,
            selectedChapterKeys: selectedChapterKeys.length > 0 && selectedChapterKeys[0] === record.id ? [] : [record.id],
        });
    }

    additionalTable = () => {
        const { selectedChapterKeys, selectedPageId, hasSelectedPage } = this.state
        const { pagesData } = this.props
        const chapterRowSelection = {
            selectedRowKeys: selectedChapterKeys,
            onChange: this.onChapterChange,
            onSelect: this.onChapterSelect,
            type: 'radio',
        }


        const chapterColumns = [
            {
                title: '文章标题',
                dataIndex: 'title',
                key: 'title',
                width: '50%',
                render: (text, record) => {
                    if (record.id === -1 || (record.id === selectedChapterKeys[0])) {
                        return <EditableCell dataIndex='chapter.title' value={record.phone} onChange={this.onNewRowChange} />
                    }
                    return text
                }
            },
            {
                title: '摘要',
                dataIndex: 'digest',
                key: 'digest',
                width: '50%',
                render: (text, record) => {
                    if (record.id === -1 || (record.id === selectedChapterKeys[0])) {
                        return <EditableCell dataIndex='chapter.digest' value={record.job} onChange={this.onNewRowChange} />
                    }
                    return text
                }
            },
        ];

        let chaptersData = []
        if (pagesData && pagesData.data && pagesData.data.templatePages) {
            for (var page in pagesData.data.templatePages) {
                if (page && page.chapter_list) {
                    if (page.id == selectedChapterKeys[0]) {
                        chaptersData = [...page.chapter_list.map(item => { item.key = item.id; return item })]
                    }
                }
            }
        }

        const hasSelectedChapter = selectedChapterKeys.length > 0 && selectedChapterKeys[0] !== -1

        return (
            <Tabs defaultActiveKey="1">
                <TabPane tab="文章" key="1">
                    <div style={{ marginBottom: 16 }}>
                        <Button type="primary" onClick={this.handleAddChapter}
                            disabled={!hasSelectedPage}
                        >新增</Button>
                        <Button type="primary" onClick={this.handleDeleteChapter}
                            disabled={!hasSelectedChapter}
                        >删除</Button>
                        <Modal
                            title="警告"
                            visible={this.state.chapterVisible}
                            onOk={this.hideModal}
                            onCancel={this.hideModal}
                            okText="确认"
                            cancelText="取消"
                        >
                            <p>确认删除</p>
                        </Modal>
                    </div>
                    <Table size="small" columns={chapterColumns} dataSource={chaptersData}
                        rowSelection={chapterRowSelection} pagination={false}
                        onRow={(record) => ({
                            onClick: () => this.onUserClick(record),
                        })}
                    />
                </TabPane>
            </Tabs>
        )
    }


    render() {
        const { loading, selectedRowKeys, selectedTown, editable, hasNewRow,
            currentPage, pageSize, total, expandedRowKeys,
            selectedPage, selectedPageId } = this.state;
        const { pagesData } = this.props
        const rowSelection = {
            selectedRowKeys,
            onChange: this.onSelectChange,
            type: 'radio',
        };

        let pagesWrappedData = []
        if (pagesData.data && pagesData.data.templatePages) {
            pagesWrappedData = [...pagesData.data.templatePages.map(item => { item.key = item.id; return item })]
        }

        if (hasNewRow) {
            pagesWrappedData = [{
                key: -1,
                id: -1,
                name: '',
                create_at: '',
            }, ...pagesWrappedData]
        } else {
            pagesWrappedData = [...pagesWrappedData.filter(item => item.id !== -1)]
        }

        const hasSelected = selectedRowKeys.length > 0 && selectedRowKeys[0] !== -1

        const pageColumns = [{
            title: '名称',
            dataIndex: 'name',
            width: "100%",
            render: (text, record) => {
                if (record.id === -1 || (editable && record.id === selectedRowKeys[0])) {
                    return <EditableCell dataIndex='templatepage.name' value={record.name} onChange={this.onNewRowChange} onCancel={this.handleCancelEditRow} />
                }
                return <a>{text}</a>
            }
        }];

        return (
            <div className="gutter-example">
                <BreadcrumbCustom first="页面模板管理" />
                <Row gutter={16}>
                    <Col className="gutter-row" md={10}>
                        <div className="gutter-box">
                            <Card title="页面模板列表" bordered={false}>
                                <div style={{ marginBottom: 16 }}>
                                    <Button type="primary" onClick={this.handleAdd}
                                    >新增</Button>
                                    <Button type="primary" onClick={this.handleModify}
                                        disabled={!hasSelected}
                                    >修改</Button>
                                    <Button type="primary" onClick={this.showModal}
                                        disabled={!hasSelected}
                                    >删除</Button>
                                    <Modal
                                        title="警告"
                                        visible={this.state.visible}
                                        onOk={this.handleDelete}
                                        onCancel={this.hideModal}
                                        okText="确认"
                                        cancelText="取消"
                                    >
                                        <p>确认删除页面模板：{selectedPage}</p>
                                    </Modal>
                                </div>
                                <Table rowSelection={rowSelection}
                                    size="small"
                                    columns={pageColumns}
                                    dataSource={pagesWrappedData}
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
                    <Col className="gutter-row" md={10}>
                        <div className="gutter-box">
                            <Card title={selectedPage ? "页面模板 " + selectedPage : "请选择页面模板"} bordered={false}
                                bodyStyle={{ paddingTop: 0 }}>
                                <div style={{}}>
                                </div>
                                {this.additionalTable()}
                            </Card>
                        </div>
                    </Col>
                </Row>
            </div>
        )
    }
}

const mapStateToProps = state => {
    const {
        pagesData = { data: { count: 0, templatePages: [] } },
    } = state.httpData;
    return { pagesData }
};

const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch)
});

export default connect(mapStateToProps, mapDispatchToProps)(PageManager)

