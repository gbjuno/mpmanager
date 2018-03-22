/**
 * Created by Jingle on 2017/11/4.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { Row, Col, Card, Button, Icon, Modal, message } from 'antd';
import * as _ from 'lodash'
import { fetchData, receiveData, updateMenu, deleteMenu } from '../../action';
import * as CONSTANTS from '../../constants';
import BreadcrumbCustom from '../../components/BreadcrumbCustom';
import * as config from '../../axios/config'
import * as utils from '../../utils'
import MenuForm from './MenuForm'

class MenuManager extends React.Component {
    state = {
        rate: 1,
        standardHeight: 200,
        baseHeight: 0,
        selectedMenuKey: null,
        selectedMenu: null,
        filter: {},
        visible: false,
    };

    componentDidMount = () => {
        this.resize();
        window.onresize = () =>{
            this.resize();
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
        const { fetchData, updateMenu } = this.props
        fetchData({funcName: 'fetchMenus', stateName: 'menusData'})
            .then(res => {
                let resData = res.data
                let i = 1;
                resData.button.map(b => {
                    b.frontend_key = i + "";
                    i++;
                    let j = 1;
                    if(b.sub_button){
                        b.sub_button.map(sb => {
                            sb.frontend_key = i + "-" + j;
                            j++;
                            return sb;
                        })
                    }
                    return b;
                })
                console.log('from wechat api---', resData)
                updateMenu(resData, null)
            })
    }


    getClientWidth = () => {    // 获取当前浏览器宽度并设置responsive管理响应式
        const { receiveData } = this.props;
        const clientWidth = document.body.clientWidth;
        receiveData({isMobile: clientWidth <= 992}, 'responsive');
    };


    resize = () => {
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

    showModal = () => {
        this.setState({
            visible: true,
        })
    }

    hideModal = () => {
        this.setState({
            visible: false,
        })
    }

    handleMenuClick = (menu, isSub, isMenuOpacity) => {
        console.log('ccccccc', menu, isSub, isMenuOpacity)
        if(isMenuOpacity) return
        if(isSub && menu.type === 'new') return
        this.setState({
            selectedMenuKey: menu.frontend_key,
            selectedMenu: menu,
        })
        if(menu.type === 'newButton' && isSub){
            this.handleAddSubMenu(menu)
        }else if(menu.type === 'newButton' && !isSub){
            this.handleAddMenu(menu)
        }
    }

    handleAddMenu = (menu) => {
        const { updateMenu, wechatLocal } = this.props
        updateMenu(wechatLocal.mergedMenus, menu, true, false)

        console.log('add main menu...', menu)

        this.setState({
            selectedMenuKey: menu.frontend_key,
            selectedMenu: menu,
        })
    }

    handleAddSubMenu = (subMenu) => {
        const { updateMenu, wechatLocal } = this.props
        updateMenu(wechatLocal.mergedMenus, subMenu, true, true)
        let sm_fk = subMenu.frontend_key
        let sm_fk_l = sm_fk.split('-')[0]
        let sm_fk_r = sm_fk.split('-')[1]
        let new_fk = sm_fk_l + "-" + (parseInt(sm_fk_r) - 1)
        subMenu.frontend_key = new_fk

        this.setState({
            selectedMenuKey: new_fk,
            selectedMenu: subMenu,
        })
    }

    handleDeleteMenu = (selectedMenu) => {
        const { deleteMenu, wechatLocal } = this.props
        console.log('selected delete menu', selectedMenu)
        deleteMenu(wechatLocal.mergedMenus, selectedMenu)
        this.setState({
            visible: false,
        })
    }

    saveAndPublish = () => {
        const { wechatLocal, fetchData } = this.props
        fetchData({funcName:'saveMenus', params: {...wechatLocal.mergedMenus},stateName:'saveMenuStatus'})
            .then(res => {
                this.fetchMenusData()
                message.info('发布成功')
            })
    }

    genSubMenus = (menu) => {
        let newSubMenus = []
        if(menu.sub_button) newSubMenus =  menu.sub_button.slice(0)
        const len = newSubMenus.length
        for(let i=len; i<4; i++){
            newSubMenus.unshift(
                {
                    "type": "new", 
                    "name": "未定义", 
                    "url": "", 
                    "sub_button": [ ]
                }
            )
        }
        if(len < 5){
            newSubMenus.push(
                {
                    "type": "newButton", 
                    "name": "子菜单名称", 
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
            return m
        })
        return newSubMenus
    }

    genMenuList = (menusData) => {
        const { selectedMenuKey, selectedMenu } = this.state
        
        console.log('fffffff', menusData)

        let buttons = [] 
        if(menusData.button){
            menusData.button.map(m => {
                if(m.sub_button && m.sub_button.length > 0){
                    delete m['url']
                }
                return m
            })
            buttons = menusData.button.slice(0)
        }
        if(buttons && buttons.length < 3){
            buttons.push(
                {
                    frontend_key: buttons.length + 1,
                    "type": "newButton", 
                    "name": "菜单名称", 
                    "sub_button": [ ],
                    "url": "",
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
                                style={{opacity: this.isMenuOpacity(subMenu, selectedMenu)?0:1}}
                                onClick={this.handleMenuClick.bind(this, subMenu, true, this.isMenuOpacity(subMenu, selectedMenu))}
                            >
                                {subMenu.type ==='newButton'?
                                <Icon style={{fontSize: 14, fontWeight: 'bold'}} type="plus" />
                                :
                                subMenu.name
                                }
                            </Row>
                            
                        ))
                    }
                    <div style={{opacity:this.isMenuOpacity(menu, selectedMenu)?0:1, height: 9}}>
                    <i className="arrow arrow_out" />
                    <i className="arrow arrow_in" />
                    </div>
                    {menu.type !== "newButton"?
                    <Row className={"wechat-main-menu " + (selectedMenuKey === menu.frontend_key?"wechat-main-menu-selected":"wechat-main-menu-unselected")} 
                        onClick={this.handleMenuClick.bind(this, menu, false, false)}
                    >
                        {menu.name}
                    </Row>
                    :
                    <Row className={"wechat-main-menu " + (selectedMenuKey === menu.frontend_key?"wechat-main-menu-selected":"wechat-main-menu-unselected")} 
                        onClick={this.handleMenuClick.bind(this, menu, false, false)}
                    >
                        <Icon style={{fontSize: 14, fontWeight: 'bold'}} type="plus" />
                    </Row>
                    }
                </Col>
            )
            })
        }else{
            return (
                <Col className="wechat-menu-row-item" md={12}><Icon style={{fontSize: 14, fontWeight: 'bold'}} type="plus" />添加菜单</Col>
            )
        }
    }

    isMenuOpacity = (menu, selectedMenu) => {
        if(menu.type === 'new'){
            return true
        }

        if(selectedMenu === null){
            return true
        }

        let menuKeyPrefix = menu.frontend_key.toString().split('-')[0]
        let selectedMenuKeyPrefix = selectedMenu.frontend_key.toString().split('-')[0]
        if(menuKeyPrefix === selectedMenuKeyPrefix){
            return false
        }
        return true
    }

    render() {
        const { baseHeight, selectedMenu } = this.state
        const { detailRecord, wechatLocal } = this.props

        console.log('wulun duome langbei douxihuanni menu', wechatLocal)
        console.log('wulun duome langbei douxihuanni selectedMenu', selectedMenu)
        let menusData = {}
        if(wechatLocal && wechatLocal.mergedMenus){
            menusData = wechatLocal.mergedMenus
        }

        let wrappedMenusData = this.genMenuList(menusData)

        const title = detailRecord?detailRecord.name:''
        let comment

        return (
            <div id="scPic" className="button-demo">
            <BreadcrumbCustom first="菜单管理" second="" />
                <Row gutter={20}>
                    <Col className="gutter-row" md={8}>
                        <Card 
                            className="comment-card"
                            bodyStyle={{}} 
                        >
                            <div style={{height: baseHeight-60}}>
                                <Row className="wechat-menu-row">
                                    {wrappedMenusData}
                                </Row>
                            </div>
                        </Card>
                    </Col>
                    <Col className="gutter-row" md={16}>
                        <Card 
                            extra={selectedMenu&&<a onClick={this.showModal}>删除菜单</a>}
                            title={selectedMenu?selectedMenu.name:'菜单'}
                        >
                            <div 
                                style={{height: baseHeight-113}}
                            >
                                <MenuForm menuFrontendKey={selectedMenu?selectedMenu.frontend_key:1} menu={selectedMenu}/>
                            </div>
                        </Card>
                        
                    </Col>
                </Row>
                <Row style={{marginTop: 10}}>
                    <Card 
                        className="comment-card"
                        bodyStyle={{}}>
                        <div style={{height: 40, textAlign: 'center'}}>
                            <Button type="primary" onClick={this.saveAndPublish}>保存并发布</Button>
                            <Button>重置</Button>
                        </div>
                    </Card>
                </Row>
                <Modal
                    title="警告"
                    visible={this.state.visible}
                    onOk={this.handleDeleteMenu.bind(this, selectedMenu)}
                    onCancel={this.hideModal}
                    okText="确认"
                    cancelText="取消"
                    >
                    <p>确认删除菜单：{selectedMenu?selectedMenu.name:''}</p>
                </Modal>
            </div>
        )
    }
}

const mapStateToProps = state => {
    return { wechatLocal: state.wechatLocal };
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch),
    updateMenu: bindActionCreators(updateMenu, dispatch),
    deleteMenu: bindActionCreators(deleteMenu, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(MenuManager);