import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { Row, Col, Input, Card, Upload, Icon, Button, message } from 'antd';
import RichEditor from '../../components/rich-editor'
import { fetchData, receiveData, searchFilter, resetFilter, handleArticleAttribute } from '../../action';
import * as config from '../../axios/config'
import BreadcrumbCustom from '../../components/BreadcrumbCustom';

const { TextArea } = Input

const getBase64 = (img, callback) => {
    const reader = new FileReader();
    reader.addEventListener('load', () => callback(reader.result));
    reader.readAsDataURL(img);
}

class ArticleForm extends React.Component {

    state = {
        htmlContent: ``,
        responseList: [],
        responseImageList: null,
        coverLoading: false,
    }

    receiveHtml = (content) => {
        const { handleArticleAttribute } = this.props
        console.log("recieved HTML content", content);

        handleArticleAttribute('content', content)
        this.setState({responseList:[]});
    }

    handleCoverChange = (info) =>{
        const { handleArticleAttribute } = this.props
        if (info.file.status === 'uploading') {
            this.setState({ coverLoading: true });
            return;
        }
        if (info.file.status === 'done') {
            message.success(`${info.file.name}上传成功`);
            
            handleArticleAttribute('thumb_media_id', info.file.response.media_id)
            handleArticleAttribute('thumb_url', info.file.response.url)

            getBase64(info.file.originFileObj, imageUrl => this.setState({
                imageUrl,
                coverLoading: false,
            }));
        } else if (info.file.status === 'error') {
            message.error(`${info.file.name}上传失败`);
        }
    }


    handleChange = (attribute, e) => {
        const { handleArticleAttribute } = this.props
        handleArticleAttribute(attribute, e.target.value)
    }

    beforeCoverUpload = (file) => {
        const isJPG = file.type === 'image/jpeg';
        if (!isJPG) {
          message.error('只能上传JPG格式的文件!');
        }
        const isLt2M = file.size / 1024 / 1024 < 2;
        if (!isLt2M) {
          message.error('上传图片大小不能超过2MB!');
        }
        return isJPG && isLt2M;
    }

    saveArticle = () => {
        const { fetchData } = this.props
        const { article } = this.props.wechatLocal
        console.log('the saving article ####', article)
        if(article){
            const { title, digest, author, content, thumb_media_id } = article
            if(!title){
                message.error('请输入标题')
                return
            }
            if(!digest){
                message.error('请输入摘要')
                return
            }
            if(!author){
                message.error('请输入作者')
                return
            }
            if(!content){
                message.error('请输入文章内容')
                return
            }
            if(!thumb_media_id){
                // message.error('请上传封面图片')
                // return
            }
            fetchData({funcName:'newArticle', params: article, stateName: 'newArticleStatus'}).then(res => {
                message.success('保存文章成功')
            })
        }
    }

    render() {
        const { fileList, imageUrl, coverLoading, responseImageList } = this.state
        const { wechatLocal } = this.props

        if(wechatLocal)console.log('ba xiao shuo xie wan', responseImageList)

        let policy = "";

        const uploadImageProps = {
            name: 'uploadImage',
            action: config.WECHAT_UPLOAD_METERIAL_IMAGE,
            onStart: (file) => {
                console.log('onStart', file.name);
                // this.refs.inner.abort(file);
            },
            onSuccess: (file) => {
                console.log('onSuccess', file);
                this.setState({
                    responseImageList: [{
                        key: file.media_id,
                        url: file.url,
                    }]
                })
            },
            onProgress: (step, file) => {
                console.log('onProgress', Math.round(step.percent), file.name);
            },
            onError: (err) => {
                console.log('onError', err);
            },
            withCredentials: true,
            listType: 'picture',
            fileList: responseImageList,
            data: (file) => {
            },
            multiple: true,
            beforeUpload: this.beforeCoverUpload,
            showUploadList: false,
        }
        const uploadVideoProps = {
            name: 'uploadVideo',
            action: config.WECHAT_UPLOAD_METERIAL_VIDEO,
            onChange: this.onChange,
            onStart: (file) => {
                console.log('onStart', file.name);
                // this.refs.inner.abort(file);
            },
            onSuccess: (file) => {
                console.log('onSuccess', file);
                this.setState({
                    responseImageList: [{
                        key: file.media_id,
                        url: file.url,
                    }]
                })
            },
            onProgress: (step, file) => {
                console.log('onProgress', Math.round(step.percent), file.name);
            },
            onError: (err) => {
                console.log('onError', err);
            },
            withCredentials: true,
            listType: 'picture',
            fileList: this.state.responseList,
            data: (file) => {

            },
            multiple: true,
            beforeUpload: this.beforeUpload,
            showUploadList: true,
        }
        const uploadAudioProps = {
            action: config.WECHAT_UPLOAD_METERIAL_IMAGE,
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
                <Icon type={coverLoading ? 'loading' : 'plus'} />
                <div className="ant-upload-text">上传封面</div>
            </div>
        );
        
        return (
            <div className="button-demo">
                <BreadcrumbCustom first="微信管理" second="文章管理" />
                <Card title="文章编辑" bordered={false}
                    extra={<span>
                        <Button onClick={this.saveArticle} type="primary">保存</Button>
                        <Button onClick={this.showModal}>取消</Button>
                        </span>}
                >
                <Row gutter={16}>
                    <Col md={18}>
                        <Input 
                            placeholder="请输入标题" 
                            className="wechat-article-ipt wechat-article-title" 
                            onChange={this.handleChange.bind(this, 'title')}
                        />
                        <Input 
                            placeholder="请输入作者" 
                            className="wechat-article-ipt"
                            onChange={this.handleChange.bind(this, 'author')}
                        />
                        <TextArea 
                            placeholder="请输入摘要" 
                            className="wechat-article-txa" 
                            autosize={{ minRows: 3, maxRows: 3 }}
                            onChange={this.handleChange.bind(this, 'digest')}
                        />
                    </Col>
                    <Col md={6}>
                    <Upload
                        name="uploadImage"
                        accept="image/*"
                        action={config.WECHAT_UPLOAD_METERIAL_IMAGE}
                        withCredentials
                        showUploadList={false}
                        listType="picture-card"
                        beforeUpload={this.beforeCoverUpload}
                        onPreview={this.handlePreview}
                        onChange={this.handleCoverChange}
                        className="wechat-article-upload-cover"
                    >
                        {imageUrl ? <img className="wechat-article-cover-img" src={imageUrl} alt="" /> : uploadButton}
                    </Upload>
                    </Col>
                </Row>
                </Card>
                <Row>
                    <Col className="gutter-row" md={24}>
                    {/* <LzEditor 
                        active={true}
                        lang="zh-CN"
                        importContent={this.state.htmlContent} 
                        cbReceiver={this.receiveHtml} 
                        uploadProps={uploadProps}
                    /> */}
                    <RichEditor 
                        active
                        lang="zh-CN"
                        importContent={this.state.htmlContent} 
                        cbReceiver={this.receiveHtml} 
                        uploadImageProps={uploadImageProps}
                        uploadVideoProps={uploadVideoProps}
                        uploadAudioProps={uploadAudioProps}
                    />
                    </Col>
                </Row>
            </div>
        );
    }
}

const mapStateToProps = state => {
    const { articlesData = {data: {}} } = state.httpData;
    const { wechatLocal = {} } = state
    return { articlesData, wechatLocal };
};

const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch),
    searchFilter: bindActionCreators(searchFilter, dispatch),
    resetFilter: bindActionCreators(resetFilter, dispatch),
    handleArticleAttribute: bindActionCreators(handleArticleAttribute, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(ArticleForm);
