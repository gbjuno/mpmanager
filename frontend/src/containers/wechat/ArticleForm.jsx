import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { Row, Col, Input, Card, Upload, Icon } from 'antd';
import LzEditor from 'react-lz-editor';
import { fetchData, receiveData, searchFilter, resetFilter } from '../../action';
import BreadcrumbCustom from '../../components/BreadcrumbCustom';

const { TextArea } = Input

class ArticleForm extends React.Component {

    state = {
        htmlContent: `<h1>Yankees, Peeking at the Red Sox, Will Soon Get an Eyeful</h1>
                <p>Whenever Girardi stole a glance, there was rarely any good news for the Yankees. While Girardi’s charges were clawing their way to a split of their four-game series against the formidable Indians, the Boston Red Sox were plowing past the rebuilding Chicago White Sox, sweeping four games at Fenway Park.</p>`,
        markdownContent: "## HEAD 2 \n markdown examples \n ``` welcome ```",
        responseList: [],
    }

    receiveHtml = (content) => {
        console.log("recieved HTML content", content);
        this.setState({responseList:[]});
    }

    render() {
        const { fileList, imageUrl } = this.state
        let policy = "";
        const uploadProps = {
            action: "http://v0.api.upyun.com/devopee",
            onChange: this.onChange,
            listType: 'picture',
            fileList: this.state.responseList,
            data: (file) => {

            },
            multiple: true,
            beforeUpload: this.beforeUpload,
            showUploadList: true,
        }

        const uploadButton = (
            <div>
                <Icon type="plus" />
                <div className="ant-upload-text">上传封面</div>
            </div>
        );
        
        return (
            <div className="button-demo">
                <BreadcrumbCustom first="微信管理" second="文章管理" />
                <Card title="文章编辑" bordered={false}>
                <Row gutter={16}>
                    <Col md={18}>
                        <Input placeholder="请输入标题" className="wechat-article-ipt wechat-article-title" />
                        <Input placeholder="请输入作者" className="wechat-article-ipt" />
                        <TextArea placeholder="请输入摘要" className="wechat-article-txa" autosize={{ minRows: 3, maxRows: 3 }} />
                    </Col>
                    <Col md={6}>
                    <Upload
                        accept="image/*"
                        action="//jsonplaceholder.typicode.com/posts/"
                        listType="picture-card"
                        onPreview={this.handlePreview}
                        onChange={this.handleChange}
                        className="wechat-article-upload-cover"
                    >
                        {imageUrl ? <img src={imageUrl} alt="" /> : uploadButton}
                    </Upload>
                    </Col>
                </Row>
                </Card>
                <Row>
                    <Col className="gutter-row" md={24}>
                    <LzEditor active={true} importContent={this.state.htmlContent} cbReceiver={this.receiveHtml} uploadProps={uploadProps}
                    lang="en" />
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

export default connect(mapStateToProps, mapDispatchToProps)(ArticleForm);
