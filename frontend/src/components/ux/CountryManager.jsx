/**
 * Created by hao.cheng on 2017/4/16.
 */
import React from 'react';
import { Table, Button, Row, Col, Card } from 'antd';
import { getPros } from '../../axios';
import BreadcrumbCustom from '../BreadcrumbCustom';

const columns = [{
    title: '镇名',
    dataIndex: 'username',
    width: 100,
    render: (text, record) => <a href={record.url} target="_blank">{text}</a>
}, {
    title: '村名',
    dataIndex: 'lang',
    width: 80
}, {
    title: 'star',
    dataIndex: 'starCount',
    width: 80
}, {
    title: '描述',
    dataIndex: 'description',
    width: 200
}];

class CountryManager extends React.Component {
    state = {
        selectedRowKeys: [],  // Check here to configure the default column
        loading: false,
        data: []
    };
    componentDidMount() {
        this.start();
    }
    start = () => {
        this.setState({ loading: true });
        getPros().then(res => {
            this.setState({
                data: [...res.data.map(val => {
                    val.key = val.id;
                    return val;
                })],
                loading: false
            });
        });
    };
    onSelectChange = (selectedRowKeys) => {
        console.log('selectedRowKeys changed: ', selectedRowKeys);
        this.setState({ selectedRowKeys });
    };
    render() {
        const { loading, selectedRowKeys } = this.state;
        const rowSelection = {
            selectedRowKeys,
            onChange: this.onSelectChange,
        };
        const hasSelected = selectedRowKeys.length > 0;
        return (
            <div className="gutter-example">
                <BreadcrumbCustom first="安监管理" second="村镇管理" />
                <Row gutter={16}>
                    <Col className="gutter-row" md={24}>
                        <div className="gutter-box">
                            <Card title="异步表格--GitHub今日热门javascript项目" bordered={false}>
                                <div style={{ marginBottom: 16 }}>
                                    <Button type="primary" onClick={this.start}
                                            disabled={loading} 
                                    >新增</Button>
                                    <Button type="primary" onClick={this.start}
                                            disabled={loading} 
                                    >修改</Button>
                                    <Button type="primary" onClick={this.start}
                                            disabled={loading} 
                                    >删除</Button>
                                    <span style={{ marginLeft: 8 }}>{hasSelected ? `Selected ${selectedRowKeys.length} items` : ''}</span>
                                </div>
                                <Table rowSelection={rowSelection} columns={columns} dataSource={this.state.data} />
                            </Card>
                        </div>
                    </Col>
                </Row>
            </div>
        );
    }
}

export default CountryManager;