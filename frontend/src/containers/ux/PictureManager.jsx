/**
 * Created by Jingle on 2017/11/4.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { Row, Col, Card, Tabs, Pagination } from 'antd';
import * as _ from 'lodash'
import moment from 'moment';
import { fetchData, receiveData } from '../../action';
import * as CONSTANTS from '../../constants';
import BreadcrumbCustom from '../../components/BreadcrumbCustom';
import PictureSearch from '../search/PictureSearch'
import * as config from '../../axios/config'
import * as utils from '../../utils'
import PhotoSwipe from 'photoswipe';
import PhotoswipeUIDefault from 'photoswipe/dist/photoswipe-ui-default';

import 'photoswipe/dist/photoswipe.css';
import 'photoswipe/dist/default-skin/default-skin.css';

const TabPane = Tabs.TabPane;
const DEFAULT_PIC_URL = '/html/static/null.png'

class PictureManager extends React.Component {
    state = {
        gallery: null,
        rate: 1,
        standardHeight: 200,
        placeTypes: [],
        placesData: [],
        selectedDay: moment(new Date()).format(CONSTANTS.DATE_QUERY_FORMAT),
        filter: {},
    };

    componentDidMount = () => {
        this.resizePicture();
        window.onresize = () =>{
            this.resizePicture();
        };

        this.setState({
            selectedDay: this.getSelectedDate(),
        })

        // TODO: 需要同时能再render方法中获取到以下两个数据
        this.fetchPlaceType();
        
        //this.fetchPictureData();
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

    componentWillUnmount = () => {
        this.closeGallery();
    };

    getClientWidth = () => {    // 获取当前浏览器宽度并设置responsive管理响应式
        const { receiveData } = this.props;
        const clientWidth = document.body.clientWidth;
        console.log(clientWidth);
        receiveData({isMobile: clientWidth <= 992}, 'responsive');
    };


    resizePicture = () => {
        this.getClientWidth();
        const scPic = document.getElementById("scPic");
        if(scPic === undefined || scPic === null) return;
        const swidth = document.body.clientWidth - 200;
        const benchmark = 1680
        this.setState({
            rate: swidth / benchmark,
        });
        
    }

    /** 查询条件组装 */
    getSelectedDate = () => {
        if(this.props.filter === undefined 
            || this.props.filter.picture === undefined
            || this.props.filter.picture.date === undefined){
            return moment(new Date()).format(CONSTANTS.DATE_QUERY_FORMAT)
        }else{
            return this.props.filter.picture.date
        }  
    }

    fetchPlaceType = () => {
        const { fetchData } = this.props
        fetchData({funcName: 'fetchPlaceTypes', stateName: 'placeTypes'}).then(res => {
            if(res === undefined || res.data === undefined || res.data.monitor_types === undefined) return
            this.setState({
                placeTypes: [...res.data.monitor_types],
            }, () => {
                this.fetchPictureData();
            })
        });
    }


    fetchPictureData = () => {
        const { fetchData  } = this.props
        const { placeTypes, filter } = this.state

        fetchData({funcName: 'fetchPicturesWithPlace', params: { day: this.state.selectedDay}, 
                stateName: 'picturesData'}).then(res => {
            if(res === undefined || res.data === undefined || res.data.monitor_places === undefined) return
            let picturesDataWithType = []
            for(let placeType of placeTypes){
                picturesDataWithType.push(
                    {
                    placeTypeId: placeType.id,
                    placeTypeName: placeType.name,
                    placesData: [...res.data.monitor_places.map(val => {
                            val.key = val.id;
                            return val;
                        }).filter(val => val.monitor_type_id === placeType.id)],
                    }
                );
            }
            this.setState({
                picturesDataWithType,
            });
        });
    }


    openGallery = (item) => {
        const items = [
            {
                src: item,
                w: 0,
                h: 0,
            }
        ];
        const pswpElement = this.pswpElement;
        const options = {index: 0};
        this.gallery = new PhotoSwipe( pswpElement, PhotoswipeUIDefault, items, options);
        this.gallery.listen('gettingData', (index, item) => {
            const _this = this;
            if (item.w < 1 || item.h < 1) { // unknown size
                var img = new Image();
                img.onload = function() { // will get size after load
                    item.w = this.width; // set image width
                    item.h = this.height; // set image height
                    _this.gallery.invalidateCurrItems(); // reinit Items
                    _this.gallery.updateSize(true); // reinit Items
                };
                img.src = item.src; // let's download image
            }
        });
        this.gallery.init();
    };
    closeGallery = () => {
        if (!this.gallery) return;
        this.gallery.close();
    };

    //转置图片数据
    transpositionToMatrix = picArray => {
        if(picArray===undefined || picArray.length === undefined || picArray.length === 0) return [[]];
        const colLen = 6
        const rowLen = Math.ceil(picArray.length / colLen)
        const mod = picArray.length % colLen
        let matrix = new Array()
        for(let i=0; i<colLen; i++){
            let uniArray = new Array()
            for(let j=0; j<rowLen; j++)
            {
                if(j * colLen + i >= picArray.length) break;
                uniArray.push(picArray[j * colLen + i])
            }
            matrix.push(uniArray)
        }
        return matrix
    };

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


    generateCard = (imgs, isMobile=false) => imgs.map(v1 => (
        v1.map(v2 => (
            <div key={v2.id} className="gutter-box" style={isMobile? {}: {height: this.state.standardHeight * this.state.rate + 80}}>
                <Card bordered={false} bodyStyle={isMobile? {padding: 0}: { padding: 0, height: this.state.standardHeight * this.state.rate + 60}}>
                    <div>
                        <img style={isMobile? {}: {height: this.state.standardHeight * this.state.rate}} 
                            onClick={() => {
                                if(this.hasPicture(v2.pictures)){
                                    return this.openGallery(config.SERVER_ROOT + this.getPicFull(v2.pictures))
                                }
                            }} 
                            alt="example" width="100%" src={config.SERVER_ROOT +  this.getPicThumb(v2.pictures)}/>
                    </div>
                    <div className="pa-s">
                        <h4 style={{marginBottom: '0em'}}>{v2.name}<span style={{paddingLeft: 5}}>{v2.monitor_place_id}</span></h4>
                        <small><a>{v2.placeName}<span style={{paddingLeft: 5}}>{v2.company_name}</span></a></small>
                    </div>
                </Card>
            </div>
        ))
    ))

    generateGrid = (datasWithType=[], isMobile) => datasWithType.map(dataWithType => {
        let imgs = this.transpositionToMatrix( dataWithType.placesData);
        const imgsTag = this.generateCard(imgs, isMobile)
        return (
        <TabPane tab={dataWithType.placeTypeName} key={dataWithType.placeTypeId}>
            <Row gutter={20}>
                <Col className="gutter-row" md={4}>
                    {imgsTag[0]}
                </Col>
                <Col className="gutter-row" md={4}>
                    {imgsTag[1]}
                </Col>
                <Col className="gutter-row" md={4}>
                    {imgsTag[2]}
                </Col>
                <Col className="gutter-row" md={4}>
                    {imgsTag[3]}
                </Col>
                <Col className="gutter-row" md={4}>
                    {imgsTag[4]}
                </Col>
                <Col className="gutter-row" md={4}>
                    {imgsTag[5]}
                </Col>
            </Row>
            <Pagination defaultCurrent={1} total={5} />
        </TabPane>
        )
    })

    /**
     * 经过混沌的洗礼，数据得以重组
     */
    chaos = (picturesData, placeTypes) => {
        if(placeTypes === undefined || picturesData === undefined 
            || _.isEmpty(placeTypes) || _.isEmpty(picturesData)
            || picturesData.data === undefined || picturesData.data.monitor_places === undefined) {
            return []
        }else{
            let picturesDataWithType =[]
            let pictures = picturesData.data
            for(let placeType of placeTypes){
                picturesDataWithType.push({
                    placeTypeId: placeType.id,
                    placeTypeName: placeType.name,
                    placesData: [...pictures.monitor_places.map(val => {
                        val.key = val.id;
                        return val;
                    }).filter(val => val.monitor_type_id === placeType.id)],
                });
            }
            return picturesDataWithType
        }
    }

    render() {
        const { rate, placeTypes } = this.state
        const { picturesData  } = this.props

        const isMobile = this.props.responsive.data.isMobile


        let chaosDataWithType = this.chaos(picturesData, placeTypes)
        let pictureGrids = this.generateGrid(chaosDataWithType, isMobile)


        return (
            <div id="scPic" className="gutter-example button-demo">
                <BreadcrumbCustom first="安监管理" second="图片管理" />
                <PictureSearch  fetchData={fetchData}/>
                <Tabs defaultActiveKey={placeTypes[0]?`${placeTypes[0].id}`:'0'}>
                {pictureGrids}
                </Tabs>
                <div className="pswp" tabIndex="-1" role="dialog" aria-hidden="true" ref={(div) => {this.pswpElement = div;} }>

                    <div className="pswp__bg" />

                    <div className="pswp__scroll-wrap">

                        <div className="pswp__container">
                            <div className="pswp__item" />
                            <div className="pswp__item" />
                            <div className="pswp__item" />
                        </div>

                        <div className="pswp__ui pswp__ui--hidden">

                            <div className="pswp__top-bar">

                                <div className="pswp__counter" />

                                <button className="pswp__button pswp__button--close" title="Close (Esc)" />

                                <button className="pswp__button pswp__button--share" title="Share" />

                                <button className="pswp__button pswp__button--fs" title="Toggle fullscreen" />

                                <button className="pswp__button pswp__button--zoom" title="Zoom in/out" />

                                <div className="pswp__preloader">
                                    <div className="pswp__preloader__icn">
                                        <div className="pswp__preloader__cut">
                                            <div className="pswp__preloader__donut" />
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <div className="pswp__share-modal pswp__share-modal--hidden pswp__single-tap">
                                <div className="pswp__share-tooltip" />
                            </div>

                            <button className="pswp__button pswp__button--arrow--left" title="Previous (arrow left)" />

                            <button className="pswp__button pswp__button--arrow--right" title="Next (arrow right)" />

                            <div className="pswp__caption">
                                <div className="pswp__caption__center" />
                            </div>

                        </div>

                    </div>

                </div>
                <style>{`
                    .ant-card-body img {
                        cursor: pointer;
                    }
                `}</style>
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

export default connect(mapStateToProps, mapDispatchToProps)(PictureManager);