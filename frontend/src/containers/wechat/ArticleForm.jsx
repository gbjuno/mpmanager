import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { Row, Col, Input, Card, Upload, Icon, Button, message } from 'antd';
import LzEditor from 'react-lz-editor';
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
        htmlContent: `<h1>Facebook面临州级诉讼 每泄漏一位用户数据就罚5万美元</h1>
                <p>腾讯科技讯 3月25日据国外媒体报道，Facebook与Cambridge Analytica在用户数据滥用泄漏事件之后，面临了大量关于非自愿数据共享的私人诉讼，但现在两家公司不得不面临一场来自于州级诉讼的新官司。

                最近，美国伊利诺伊州库克县已经对两家公司提起诉讼，指控其违反了该州的《消费者欺诈与欺骗性商业行为法》。在起诉书中指出，Cambridge Analytica违反了当地法律，将一款名叫“这是你的数字生活”的学术研究App作为收集用户个人数据的工具，同时违反了Facebook的隐私协议。而Facebook则被指控在了解Cambridge Analytica行为后使用错误的方式保护用户数据，同时没有采取任何措施阻止Cambridge Analytica的行为。
                
                库克县并未要求具体的赔偿总额，但根据该州法律规定，每一次违反伊利诺伊州欺诈法都会被处以5万美元的罚款。如果受害者的年龄超过65周岁，还要被追加1万美元。虽然目前不确定伊利诺伊州究竟有多少人在这5000万泄露的个人用户数据名单之列，但如果诉讼成功的话，对Facebook和Cambridge Analytica公司来说，都是非常高昂的代价。
                
                当Facebook首席执行官马克·扎克伯格和首席运营官雪莉·桑德伯格承认公司在处理Cambridge Analytica分析工具的行为出现失误时，表示社交网络能做的就是尽量减少损害而无法做到完全避免损害。但如果有其它州也加入到对两家公司的诉讼行为中来也不足为奇。毕竟此次数据泄露事件涉及5000万用户，除了伊利诺伊州之外一定会有来自于其它地区的更多受害者</p>`,
        markdownContent: "## HEAD 2 \n markdown examples \n ``` welcome ```",
        responseList: [],
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
                message.error('请上传封面图片')
                return
            }
            fetchData({funcName:'newArticle', params: article, stateName: 'newArticleStatus'}).then(res => {
                message.success('保存文章成功')
            })
        }
    }

    render() {
        const { fileList, imageUrl, coverLoading } = this.state
        const { wechatLocal } = this.props

        if(wechatLocal)console.log('ba xiao shuo xie wan', wechatLocal.article)

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
                        withCredentials={true}
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
                    <LzEditor 
                        active={true}
                        lang="zh-CN"
                        importContent={this.state.htmlContent} 
                        cbReceiver={this.receiveHtml} 
                        uploadProps={uploadProps}
                    lang="en" />
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
