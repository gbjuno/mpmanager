/**
 * Created by Jingle on 2017/11/4.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { Row, Col, Card } from 'antd';
import { fetchData, receiveData } from '../../action';
import BreadcrumbCustom from '../BreadcrumbCustom';
import PictureSearch from './search/PictureSearch'
import * as config from '../../axios/config'
import PhotoSwipe from 'photoswipe';
import PhotoswipeUIDefault from 'photoswipe/dist/photoswipe-ui-default';

import 'photoswipe/dist/photoswipe.css';
import 'photoswipe/dist/default-skin/default-skin.css';

class PictureManager extends React.Component {
    state = {
        gallery: null,
        rate: 1,
        responsive: false,
        picturesData: [],
        picturesDataWithType: [],
        standardHeight: 430,
    };
    componentDidMount = () => {
        this.resizePicture();
        const clientWidth = document.body.clientWidth;
        if(clientWidth <= 992) {
            this.setState({
                responsive: true,
            })
        }
        window.onresize = () =>{
            const clientWidth = document.body.clientWidth;
            if(clientWidth <= 992) {
                this.setState({
                    responsive: true,
                })
                return;
            }else{
                this.resizePicture();
            }
            
        };

       this.fetchPlaceType();
    };

    fetchPlaceType = () => {
        const { fetchData } = this.props
        let tempTownId
        fetchData({funcName: 'fetchPlaceTypes', stateName: 'placeTypes'}).then(res => {
            if(res === undefined || res.data === undefined || res.data.monitor_types === undefined) return
            this.setState({
                placeTypes: [...res.data.monitor_types.map(val => {
                    val.key = val.id;
                    return val;
                })],
                placeTypeLoading: false,
            }, () => {
                this.fetchPictureData();
            });
        });
    }

    fetchPictureData = () => {
        const { fetchData } = this.props
        const { placeTypeLoading, placeTypes } = this.state
        fetchData({funcName: 'fetchPictures', stateName: 'picturesData'}).then(res => {
            if(res === undefined || res.data === undefined || res.data.pictures === undefined) return
            let picturesDataWithType = []
            for(let placeType of placeTypes){
                console.log('placeType.id', placeType.id)
                picturesDataWithType.push(
                    {
                    placeTypeId: placeType.id,
                    placeTypeName: placeType.name,
                    picturesData: [...res.data.pictures.map(val => {
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

    componentDidUpdate = (nextProps, nextState) => {
    };

    componentWillUnmount = () => {
        this.closeGallery();
    };


    resizePicture = () => {
        const scPic = document.getElementById("scPic");
        if(scPic === undefined || scPic === null) return;
        const swidth = scPic.clientWidth;
        const benchmark = 1680;
        this.setState({
            rate: swidth / benchmark,
            responsive: false,
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

    generateCard = imgs => imgs.map(v1 => (
        v1.map(v2 => (
            <div key={v2.id} className="gutter-box" style={this.state.responsive? {}: {height: this.state.standardHeight * this.state.rate + 80}}>
                <Card bordered={false} bodyStyle={this.state.responsive? {padding: 0}: { padding: 0, height: this.state.standardHeight * this.state.rate + 60}}>
                    <div>
                        <img style={this.state.responsive? {}: {height: this.state.standardHeight * this.state.rate}} onClick={() => this.openGallery(config.SERVER_ROOT + v2.full_uri)} 
                            alt="example" width="100%" src={config.SERVER_ROOT + v2.full_uri} />
                    </div>
                    <div className="pa-m">
                        <h3>{v2.companyName}<span style={{paddingLeft: 5}}>{v2.monitor_place_id}</span></h3>
                        <small><a>{v2.placeName}<span style={{paddingLeft: 5}}>{v2.create_at.substring(0, 10)}</span></a></small>
                    </div>
                </Card>
            </div>
        ))
    ))

    generateGrid = datasWithType => datasWithType.map(dataWithType => {
        let imgs = this.transpositionToMatrix( dataWithType.picturesData);
        const imgsTag = this.generateCard(imgs)
        return (
        <div key={dataWithType.placeTypeId}>
            <h2>{dataWithType.placeTypeName}</h2>
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
        </div>
        )
    })

    render() {
        const { rate, responsive, picturesDataWithType } = this.state
        const { picData, fetchData} = this.props
        let pictureGrids = this.generateGrid(picturesDataWithType)
        console.log('pictureGrids ---> ', pictureGrids)
        
        return (
            <div id="scPic" className="gutter-example button-demo">
                <BreadcrumbCustom first="安监管理" second="图片管理" />
                <PictureSearch style={{paddingBottom: 13}} fetchData={fetchData}/>
                {pictureGrids}
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
    const { picData = {data: []} } = state.httpData;
    return { picData };
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch)
});

export default connect(mapStateToProps, mapDispatchToProps)(PictureManager);