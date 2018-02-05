/**
 * Created by Haiyang on 2018/2/4.
 */
import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import BannerAnim, { Element, Thumb } from 'rc-banner-anim';
import * as _ from 'lodash'
import moment from 'moment';
import { fetchData, receiveData, searchFilter } from '../../action';
import TweenOne from 'rc-tween-one';
import * as config from '../../axios/config'

import PhotoSwipe from 'photoswipe';
import PhotoswipeUIDefault from 'photoswipe/dist/photoswipe-ui-default';

const BgElement = Element.BgElement;

class PictureCarousel extends React.Component {
    state = {
    };

    componentDidMount = () => {
        const { searchFilter, elements } = this.props
        const pictures = elements.pictures
        if(pictures === undefined || pictures.length === 0) return
        searchFilter('pictureCarousel', {selectPicture: pictures[0]})
    }

    componentWillUnmount = () => {
        this.closeGallery();
    };

    handleChange = (before, index) => {
        const { searchFilter, elements } = this.props
        const pictures = elements.pictures
        if(pictures === undefined || pictures.length === 0) return
        searchFilter('pictureCarousel', {selectPicture: pictures[index]})
    }

    onMouseEnter = () => {
        this.setState({
            enter: true,
        });
    }
    
    onMouseLeave = () => {
        this.setState({
            enter: false,
        });
    }

    generateElement = () => {
        const { elements } = this.props
        const pictures = elements.pictures
        if(pictures === undefined || pictures.length === 0) return
        return (
            pictures.map(picture => (
                <Element key={`${picture.id}`}
                    prefixCls="banner-user-elem"
                >
                    <BgElement
                    key="bg"
                    className="bg"
                    style={{
                        backgroundImage: `url(${config.SERVER_ROOT + picture.full_uri})`,
                        backgroundSize: 'auto 100%',
                        backgroundPosition: 'center',
                        backgroundRepeat: 'no-repeat',
                    }}
                    onClick={() => this.openGallery(config.SERVER_ROOT + picture.full_uri)} 
                    />
                    <TweenOne
                    animation={{ y: 50, opacity: 0, type: 'from', delay: 200 }}
                    key="2"
                    name="TweenOne"
                    >
                    </TweenOne>
                </Element>
                )
            )
        )
    }

    thumbChildren = () => { 
        const { elements } = this.props
        const pictures = elements.pictures
        if(pictures === undefined || pictures.length === 0) return
        return (pictures.map((picture, i) =>
            <span key={i}><i style={{ backgroundImage: `url(${config.SERVER_ROOT + picture.full_uri})` }} /></span>
            )
        )
    };

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

    render(){
        const { height } = this.props
        return (
            <div>
            <BannerAnim prefixCls="banner-user" style={{height: height}} onChange={this.handleChange}
                onMouseEnter={this.onMouseEnter} onMouseLeave={this.onMouseLeave}>
                {this.generateElement()}
                <Thumb prefixCls="user-thumb" key="thumb" component={TweenOne}
                animation={{ bottom: this.state.enter ? 0 : -70 }}
                >
                {this.thumbChildren()}
                </Thumb>
            </BannerAnim>
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
        );
    }
}
const mapStateToProps = state => {
    const { searchFilter } = state;
    return { filter: searchFilter };
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch),
    searchFilter: bindActionCreators(searchFilter, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(PictureCarousel);