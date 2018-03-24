/**
 * Created by Jingle Chen on 2017/12/7.
 */
import React from 'react';
import { Link } from 'react-router-dom'
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import * as _ from 'lodash'
import {  Button, Row, Col, Card, Avatar, List, Input, Divider, Icon, message } from 'antd';
import { fetchData, receiveData, searchFilter, resetFilter } from '../../action';
import { getPros } from '../../axios';
import BreadcrumbCustom from '../../components/BreadcrumbCustom';
import EditableCell from '../../components/cells/EditableCell';

const IconText = ({ type, text }) => (
    <span>
      <Icon type={type} style={{ marginRight: 8 }} />
      {text}
    </span>
  );

class ArticleManager extends React.Component {
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
        resetFilter('article')
        searchFilter('article', {
            pageNo: currentPage,
            pageSize: pageSize,
        })
        this.setState({ loading: true }, () => {
            this.fetchData();
        });
        
    };

    fetchData = () => {
        const { fetchData, searchFilter, filter } = this.props

        fetchData({funcName: 'fetchArticles', params: filter.user, stateName: 'articlesData'}).then(res => {
            console.log('from api article data...', res)
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

        const { loading, selectedRowKeys, hasNewRow, pageSize,
             editable } = this.state;
        const { articlesData, filter } = this.props;

        let total = 0; 
        let currentPage = 1;
        if(filter.user) {
            total = filter.user.total
            currentPage = filter.user.pageNo
        }
        console.log('sally, i will take all my life to protect you...', articlesData)

        let options = [];

        let articlesWrappedData = []
        if(articlesData.data && articlesData.data.chapters){
            articlesWrappedData = [...articlesData.data.chapters.map(item => {item.key = item.id; return item})]
        }

        const hasSelected = selectedRowKeys.length > 0 && selectedRowKeys[0] !== -1

        
        return (
            <div className="gutter-example">
                <BreadcrumbCustom first="微信管理" second="文章管理" />
                <Row gutter={16}>
                    <Col className="gutter-row" md={24}>
                        <div className="gutter-box">
                            <Card title="文章列表" bordered={false}>
                                <Link to="/app/wechat/wzbj" ><Button type="primary">新增</Button></Link>
                                <Divider />
                                <List
                                    itemLayout="vertical"
                                    size="large"
                                    pagination={{
                                        hideOnSinglePage: true,
                                        onChange: this.handlePageChange,
                                        current: currentPage,
                                        defaultCurrent: 1,
                                        pageSize,
                                        total,
                                    }}
                                    dataSource={articlesWrappedData}
                                    renderItem={item => (
                                    <List.Item
                                        key={item.title}
                                        actions={[<IconText type="edit" text="编辑" />, <IconText type="delete" text="删除" />, <IconText type="message" text="2" />]}
                                        extra={<img width={272} alt="logo" src={item.url} />}
                                    >
                                        <List.Item.Meta
                                        avatar={<Avatar src={item.avatar} />}
                                        title={<a href={item.href}>{item.title}</a>}
                                        description={item.digest}
                                        />
                                        {item.content}
                                    </List.Item>
                                    )}
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
        articlesData = {data: {}}, 
    } = state.httpData;
    return { articlesData, filter: state.searchFilter };
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch),
    searchFilter: bindActionCreators(searchFilter, dispatch),
    resetFilter: bindActionCreators(resetFilter, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(ArticleManager)

