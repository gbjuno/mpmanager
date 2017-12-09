/**
 * Created by Jingle on 2017/11/4.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { Row, Col, Card } from 'antd';
import { fetchData, receiveData } from '../../action';
import BreadcrumbCustom from '../BreadcrumbCustom';
import SearchForm from '../forms/SearchForm'
import PhotoSwipe from 'photoswipe';
import PhotoswipeUIDefault from 'photoswipe/dist/photoswipe-ui-default';

import 'photoswipe/dist/photoswipe.css';
import 'photoswipe/dist/default-skin/default-skin.css';

class PictureManager extends React.Component {
    state = {
        gallery: null,
        rate: 1,
        responsive: false,
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

        const { fetchData } = this.props
        fetchData({funcName: 'fetchScPic', stateName: 'picData'});
    };

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

    render() {
        const { rate, responsive } = this.state
        const { picData, fetchData} = this.props
        const imgs = this.transpositionToMatrix( picData.data);
        const standardHeight = 430
        const imgsTag = imgs.map(v1 => (
            v1.map(v2 => (
                <div key={v2.id} className="gutter-box" style={responsive? {}: {height: standardHeight * rate + 80}}>
                    <Card bordered={false} bodyStyle={responsive? {padding: 0}: { padding: 0, height: standardHeight * rate + 60}}>
                        <div>
                            <img style={responsive? {}: {height: standardHeight * rate}} onClick={() => this.openGallery(v2.src)} alt="example" width="100%" src={v2.src} />
                        </div>
                        <div className="pa-m">
                            <h3>{v2.companyName}<span style={{paddingLeft: 5}}>{v2.name}</span></h3>
                            <small><a>{v2.placeName}<span style={{paddingLeft: 5}}>{v2.createAt}</span></a></small>
                        </div>
                    </Card>
                </div>
            ))
        ));
        return (
            <div id="scPic" className="gutter-example button-demo">
                <BreadcrumbCustom first="安监管理" second="图片管理" />
                <SearchForm style={{paddingBottom: 13}} fetchData={fetchData}/>
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