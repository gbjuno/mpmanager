/**
 * Created by Jingle on 2017/11/4.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { Row, Col, Card } from 'antd';
import BreadcrumbCustom from '../BreadcrumbCustom';
import SearchForm from '../forms/SearchForm'
import PhotoSwipe from 'photoswipe';
import PhotoswipeUIDefault from 'photoswipe/dist/photoswipe-ui-default';
import { fetchData, receiveData } from '../../action';

import 'photoswipe/dist/photoswipe.css';
import 'photoswipe/dist/default-skin/default-skin.css';

class SecurityCheckPic extends React.Component {
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
        // const imgs2 = [
        //     [
        //         {src: 'http://img.hb.aicdn.com/1cad414972c5db2b8c1942289e3aeef37175006a8bb16-CBtjtX_fw', id: 'r1_c1', name: '照相机' },
        //         {src: 'http://img.hb.aicdn.com/016f2e13934397e17c3482a4529f3da1149d37fd2a99c-RVM1Gi_fw', id: 'r2_c1', name: '火影'},
        //         {src: 'http://img.hb.aicdn.com/8c5d5f2bf6427d1b5ed8657a7ae0c9938d3465e367899-AJ0zVA_fw', id: 'r3_c1', name: '夜猫'},
        //         {src: 'http://img.hb.aicdn.com/bd71ccac0b16bbcade255a1a8a63504d71c7dee9a8652-zBCN9d_fw', id: 'r4_c1', name: '蛋糕'},
        //         {src: 'http://img.hb.aicdn.com/37a40cb04345463858d45418ae6ed9ef319e30dc37a45-o4pQ0j_fw', id: 'r5_c1', name: '峡谷'},
        //     ],
        //     [
        //         {src:'http://img.hb.aicdn.com/5fad6c3a14a9b80c4448835bb6b23ab895d18e234eff3-BPGmox_fw', id: 'r1_c2', name: '美女背影与直升机'},
        //         {src:'http://img.hb.aicdn.com/a1a19de5dac212a646ba6967ef565786399fb1665bd04-EEvwzR_fw', id: 'r2_c2', name: '呆猫'},
        //         {src:'http://img.hb.aicdn.com/06595f8044e881de3a82d691768bc8c21a2a9f3633d60-XKjC2s_fw', id: 'r3_c2', name: '岸边公路'},
        //         {src:'http://img.hb.aicdn.com/880787b36d45efbe05aa409c867db29a3028e02da7f9b-qxGib9_fw', id: 'r4_c2', name: '下雨的街道'},
        //         {src:'http://img.hb.aicdn.com/4964b97f6f6eb61a20922b40842adf0169c44e491c4b60-azX1S7_fw', id: 'r5_c2', name: '山间列车'},
        //     ], 
        //     [   
        //         {src:'http://img.hb.aicdn.com/ff97d00944edfc706c62dd5c0e955c4099a37b407534f-BcUqf0_fw', id: 'r1_c3', name: '日本街道'},
        //         {src:'http://img.hb.aicdn.com/0e22be22b08c6f78b94283b6cfa890093ac3cae8401e7-b1ftfi_fw', id: 'r2_c3', name: '都市街道'},
        //         {src:'http://img.hb.aicdn.com/879f870e15f7cc0847c8ae19a5fcbe974d5904bb181d7-RGmtNU_fw', id: 'r3_c3', name: '旷野公路'},
        //         {src:'http://img.hb.aicdn.com/b4a8e62958555a97dc3de9ccb03284bf556c042925522-x50qGv_fw', id: 'r4_c3', name: '水城'},
        //         {src:'http://img.hb.aicdn.com/1ef493a15674e9fd523b248ea4ec43d2ea9ce6952ff3e-WavWKc_fw', id: 'r5_c3', name: '休闲步道'},
        //     ],  
        //     [    
        //         {src:'http://img.hb.aicdn.com/8e16efec78ac4a3684fc8999d18e3661af40fd4510a25-DDvQON_fw', id: 'r1_c4', name: '海浪'},
        //         {src:'http://img.hb.aicdn.com/61dfa024c8040e6a5bcb03d42928fbcb0c87c1a54e731-yc4lvV_fw', id: 'r2_c4', name: '大海日落'},
        //         {src:'http://img.hb.aicdn.com/6783b4d7811ad7fb87b1446c5488b91179f7608118289-hpEyP3_fw', id: 'r3_c4', name: '浪花'},
        //         {src:'http://img.hb.aicdn.com/7be61ba6bdb20a73be63edc387b16eec72d0bbb51c7ef-XafA07_fw', id: 'r4_c4', name: '海沟'},
        //         {src:'http://img.hb.aicdn.com/bd3ba3f907fe098b911947e0020615b50fc340ed2df72-WsuHuM_fw', id: 'r5_c4', name: '海岸线'},
        //     ],  
        //     [    
        //         {src:'http://img.hb.aicdn.com/dd4dd2d5d1fde0e05204a257be7710a2190c48a31c950-2NwZPC_fw658', id: 'r1_c5', name: '沙滩'},
        //         {src:'http://img.hb.aicdn.com/cb16c68c4d3b7a08b5e91cd351f6b723634ca3fc27d4d-m1JD8z_fw', id: 'r2_c5', name: '滑板'},
        //         {src:'http://img.hb.aicdn.com/e3559b6e8d7237857382050e5659a64cc0b7d696a2869-stcRXA_fw', id: 'r3_c5', name: '金门大桥'},
        //         {src:'http://img.hb.aicdn.com/4ea229436fcf2077502953907a6afb16d3c5cd611b8e2-0dVIeH_fw', id: 'r4_c5', name: '转角遇到店'},
        //         {src:'http://img.hb.aicdn.com/98c786f4314736f95a42bf927bf65a82d305a532c6258-njI6id_fw', id: 'r5_c5', name: '转角路标'},
        //     ],  
        //     [   
        //         {src:'http://img.hb.aicdn.com/73edcfbcf352f69ea07d04ae615040fa4479846f1e2c4-cVtzRr_fw658', id: 'r1_c6', name: '油画'},
        //         {src:'http://img.hb.aicdn.com/d7b555c88e331a11bc42a29cd6a8d89d3fd6d1cf2d018-r0n85r_fw658', id: 'r2_c6', name: '池田依来沙'},
        //         {src:'http://img.hb.aicdn.com/aa30d771f8515112b03a586c1ba8a87e05704f588ea27-nwgkLh_fw658', id: 'r3_c6', name: '解码语者'},
        //         {src:'http://img.hb.aicdn.com/e6d618e401de5cf490ded4112df2122496ebdb9c228a7-4VhX4S_fw658', id: 'r4_c6', name: '小熊猫'},
        //         {src:'http://img.hb.aicdn.com/dca1994d81a0120f4c00b0e0298d3fa4aa73c439fc347-9JAHpF_fw658', id: 'r5_c6', name: '二次元'},
        //     ]
        // ];
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
                <BreadcrumbCustom first="UI" second="画廊(图片来自花瓣网，仅学习，若侵权请联系删除)" />
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

export default connect(mapStateToProps, mapDispatchToProps)(SecurityCheckPic);