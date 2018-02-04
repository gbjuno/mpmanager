/**
 * Created by Jingle on 2017/11/4.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { Row, Col, Card, Button, Alert } from 'antd';
import * as _ from 'lodash'
import moment from 'moment';
import { fetchData, receiveData } from '../../action';
import * as CONSTANTS from '../../constants';
import BreadcrumbCustom from '../../components/BreadcrumbCustom';
import Carousel from '../../components/carousel';
import PictureSearch from '../search/PictureSearch'
import * as config from '../../axios/config'
import * as utils from '../../utils'

import 'photoswipe/dist/photoswipe.css';
import 'photoswipe/dist/default-skin/default-skin.css';

const DEFAULT_PIC_URL = '/html/static/null.png'

class PictureDetail extends React.Component {
    state = {
        gallery: null,
        rate: 1,
        standardHeight: 200,
        baseHeight: 0,
        selectedDay: moment(new Date()).format(CONSTANTS.DATE_QUERY_FORMAT),
        filter: {},
    };

    componentDidMount = () => {
        this.resizePicture();
        window.onresize = () =>{
            this.resizePicture();
        };

    };

    componentDidUpdate(prevProps, prevState){
        const oldFilter = prevProps.filter
        const newFilter = this.props.filter

        if( oldFilter !== newFilter ){
            this.setState({
                filter: newFilter,
            })
        }
    }

    
    componentDidUpdate = (nextProps, nextState) => {
    };


    getClientWidth = () => {    // 获取当前浏览器宽度并设置responsive管理响应式
        const { receiveData } = this.props;
        const clientWidth = document.body.clientWidth;
        receiveData({isMobile: clientWidth <= 992}, 'responsive');
    };


    resizePicture = () => {
        this.getClientWidth();
        const scPic = document.getElementById("scPic");
        if(scPic === undefined || scPic === null) return;
        const sWidth = document.body.clientWidth - 200;
        const sHeight = document.body.clientHeight;
        const benchmark = 1680
        this.setState({
            baseHeight: sHeight - 213,
            rate: sWidth / benchmark,
        });
        
    }



    hasPicture = pictures => {
        if (pictures === undefined ) return false
        if (pictures.length === 0) return false
        if (pictures[0].full_uri === undefined) return false
        if (pictures[0].thumb_uri === undefined) return false
        return true
    }

    getPicThumb = pictures => {
        let picThumb = this.hasPicture(pictures)? pictures[0].full_uri : DEFAULT_PIC_URL
        return picThumb
    }

    getPicFull = pictures => {
        let picFull = this.hasPicture(pictures)? pictures[0].full_uri : DEFAULT_PIC_URL
        return picFull
    }


    getComment = () => {
        const { detailRecord } = this.props
        if(detailRecord.pictures && detailRecord.pictures[0]){
            return detailRecord.pictures[0].judgecomment
        }
    }

    isUnqualified = () => {
        const { detailRecord } = this.props
        if(detailRecord.pictures && detailRecord.pictures[0]){
            return detailRecord.pictures[0].judgement === 'F'
        }
        return false
    }

    handleBack = () => {
        if(this.props.onBack){
            this.props.onBack()
        }
    }

    render() {
        const { rate, placeTypes, baseHeight } = this.state
        const { detailRecord  } = this.props

        const isMobile = this.props.responsive.data.isMobile

        const title = detailRecord.name
        const isUnqualified = this.isUnqualified()
        let comment 
        if(isUnqualified){
            comment = this.getComment()
        }else{
            comment = '当前为最新所拍照片'
        }

        return (
            <div id="scPic" className="gutter-example button-demo">
                <Row gutter={20}>
                    <Col className="gutter-row" md={18}>
                        <Card >
                            {baseHeight !== 0 &&
                            <Carousel elements={detailRecord} height={baseHeight}/>
                            }
                        </Card>
                    </Col>
                    <Col className="gutter-row" md={6}>
                        <Card 
                            title={title}
                            extra={<Button onClick={this.handleBack} style={{cursor:'pointer'}}>返回</Button>}
                            bodyStyle={{}}>
                            <Alert
                                message={comment}
                                type={isUnqualified?"error":"success"}
                                showIcon
                            />
                        </Card>
                    </Col>
                </Row>
            </div>
        )
    }
}

const mapStateToProps = state => {
    return { ...state.httpData, filter: state.searchFilter };
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch)
});

export default connect(mapStateToProps, mapDispatchToProps)(PictureDetail);