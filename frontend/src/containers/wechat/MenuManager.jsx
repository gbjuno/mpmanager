/**
 * Created by Jingle on 2017/11/4.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { Row, Col, Card, Button, Alert, Icon, Popover } from 'antd';
import * as _ from 'lodash'
import moment from 'moment';
import { fetchData, receiveData, searchFilter } from '../../action';
import * as CONSTANTS from '../../constants';
import BreadcrumbCustom from '../../components/BreadcrumbCustom';
import * as config from '../../axios/config'
import * as utils from '../../utils'
import MenuForm from './MenuForm'

import 'photoswipe/dist/photoswipe.css';
import 'photoswipe/dist/default-skin/default-skin.css';

const DEFAULT_PIC_URL = '/html/static/null.png'

class MenuManager extends React.Component {
    state = {
        rate: 1,
        standardHeight: 200,
        baseHeight: 0,
        selectedMenuKey: null,
        selectedMenu: null,
        filter: {},
    };

    componentDidMount = () => {
        this.resizePicture();
        window.onresize = () =>{
            this.resizePicture();
        };

        this.fetchMenusData()
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


    fetchMenusData = () => {
        const { fetchData } = this.props
        fetchData({funcName: 'fetchMenus', stateName: 'menusData'})
    }


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

    handleMenuClick = (menu, isSub) => {
        if(isSub && menu.type === 'new') return
        this.setState({
            selectedMenuKey: menu.frontend_key,
            selectedMenu: menu,
        })
    }

    genSubMenus = (menu) => {
        let newSubMenus = menu.sub_button.slice(0)
        const len = newSubMenus.length
        for(let i=len; i<5; i++){
            newSubMenus.unshift(
                {
                    "type": "new", 
                    "name": "未定义", 
                    "url": "", 
                    "sub_button": [ ]
                }
            )
        }
        let j = 0;
        const k = menu.frontend_key
        newSubMenus.map(m => {
            m.frontend_key = k + "-" + j;
            j++;
        })
        return newSubMenus
    }

    genMenuList = (menusData0) => {
        const { selectedMenuKey } = this.state
        let menusData = {
            "menu": {
                "button": [
                    {
                        "type": "click", 
                        "name": "今日歌曲", 
                        "key": "V1001_TODAY_MUSIC", 
                        "sub_button": [
                            {
                                "type": "view", 
                                "name": "搜索", 
                                "url": "http://www.soso.com/", 
                                "sub_button": [ ]
                            }
                         ]
                    }, 
                   
                    {
                        "name": "菜单", 
                        "sub_button": [
                            {
                                "type": "view", 
                                "name": "搜索", 
                                "url": "http://www.soso.com/", 
                                "sub_button": [ ]
                            }, 
                            {
                                "type": "view", 
                                "name": "视频", 
                                "url": "http://v.qq.com/", 
                                "sub_button": [ ]
                            }, 
                            {
                                "type": "click", 
                                "name": "赞一下我们", 
                                "key": "V1001_GOOD", 
                                "sub_button": [ ]
                            }
                        ]
                    }
                ]
            }
        }

        let buttons = menusData.menu.button
        if(buttons.length < 3){
            buttons.push(
                {
                    frontend_key: buttons.length + 1,
                    "type": "new", 
                    "name": "菜单名称", 
                    "sub_button": [ ]
                }
            )
        }

        let j = 0
        if(buttons && buttons.length > 0){
            return buttons.map(menu => {
                j++
                menu.frontend_key = j
                return (
                <Col key={menu.frontend_key} className="wechat-menu-row-item" md={24/buttons.length} >
                    {
                        this.genSubMenus(menu).map(subMenu => (
                            <Row key={subMenu.frontend_key} 
                                className={"wechat-sub-menu " + (selectedMenuKey === subMenu.frontend_key?"wechat-sub-menu-selected":"wechat-sub-menu-unselected") }
                                style={{opacity:subMenu.type === "new"?0:1}}
                                onClick={this.handleMenuClick.bind(this, subMenu, true)}
                                >
                                {subMenu.name}
                            </Row>
                            
                        ))
                    }
                    <div style={{opacity:menu.sub_button.length > 0?1:0, height: 9}}>
                    <i className="arrow arrow_out"></i>
                    <i className="arrow arrow_in"></i>
                    </div>
                    {menu.type !== "new"?
                    <Row className={"wechat-main-menu " + (selectedMenuKey === menu.frontend_key?"wechat-main-menu-selected":"wechat-main-menu-unselected")} 
                        onClick={this.handleMenuClick.bind(this, menu)}
                    >
                        {menu.name}
                    </Row>
                    :
                    <Row className={"wechat-main-menu " + (selectedMenuKey === menu.frontend_key?"wechat-main-menu-selected":"wechat-main-menu-unselected")} 
                        onClick={this.handleMenuClick.bind(this, menu, false)}
                    >
                        <Icon style={{fontSize: 14, fontWeight: 'bold'}} type="plus"/>
                    </Row>
                    }
                </Col>
            )
            })
        }else{
            return (
                <Col className="wechat-menu-row-item" md={12}><Icon style={{fontSize: 14, fontWeight: 'bold'}} type="plus"/>添加菜单</Col>
            )
        }
    }

    isMenuOpacity = () => {
        
    }

    render() {
        const { rate, baseHeight, selectedMenu } = this.state
        const { detailRecord, menusData  } = this.props

        console.log('wulun duome langbei douxihuanni', selectedMenu)

        let wrappedMenusData = this.genMenuList()

        const isMobile = this.props.responsive.data.isMobile

        const title = detailRecord?detailRecord.name:''
        let comment

        return (
            <div id="scPic" className="button-demo">
            <BreadcrumbCustom first="菜单管理" second="" />
                <Row gutter={20}>
                    <Col className="gutter-row" md={8}>
                        <Card 
                            className="comment-card"
                            bodyStyle={{}}>
                            <div style={{height: baseHeight-60}}>
                                <Row className="wechat-menu-row">
                                    {wrappedMenusData}
                                </Row>
                            </div>
                        </Card>
                    </Col>
                    <Col className="gutter-row" md={16}>
                        <Card 
                            title={"菜单"}>
                            <div 
                                style={{height: baseHeight-113}}
                            >
                                <MenuForm menuFrontendKey={selectedMenu?selectedMenu.frontend_key:0} menu={selectedMenu}/>
                            </div>
                        </Card>
                    </Col>
                </Row>
                <Row style={{marginTop: 10}}>
                    <Card 
                        className="comment-card"
                        bodyStyle={{}}>
                        <div style={{height: 40, textAlign: 'center'}}>
                            <Button type="primary">保存并发布</Button>
                            <Button>取消</Button>
                        </div>
                    </Card>
                </Row>
            </div>
        )
    }
}

const mapStateToProps = state => {
    console.log('you are my angel...', state)
    return { ...state.httpData, filter: state.searchFilter };
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch),
    searchFilter: bindActionCreators(searchFilter, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(MenuManager);