/**
 * Created by Jingle Chen on 2018/12/31.
 */
import React, { Component } from 'react';
import { Menu, Icon, Layout, Badge, Popover } from 'antd';
import screenfull from 'screenfull';
import * as _ from 'lodash'
import { queryString } from '../utils';
import * as config from '../axios/config'
import avater from '../style/imgs/b1.png';
import SiderCustom from './SiderCustom';
import { connect } from 'react-redux';
import { withRouter, Redirect } from 'react-router-dom';
import { bindActionCreators } from 'redux';
import { fetchData, receiveData } from '@/action';
const { Header } = Layout;
const SubMenu = Menu.SubMenu;
const MenuItemGroup = Menu.ItemGroup;

class HeaderCustom extends Component {
    state = {
        user: '',
        visible: false,
    };
    componentDidMount() {
        const QueryString = queryString();
    };
    screenFull = () => {
        if (screenfull.enabled) {
            screenfull.request();
        }

    };
    menuClick = e => {
        e.key === 'logout' && this.logout();
    };
    logout = () => {  
        const { fetchData } = this.props
        const user = localStorage.getItem('user')

        fetchData({funcName: 'authLogout', params: {'Phone': user}, stateName: 'authStatus'})
                .then(res => {
                    console.log('res--->', res)
                    if(res.data.status === 200){
                        // this.setState({
                        //     redirectToLogin: true,
                        // })
                        console.log('退出成功')
                        window.location.href = config.PAGE_CONTEXT
                        localStorage.removeItem('user');
                    }
                }).catch(err => {
                    let errRes = err.response
                    if(errRes && errRes.data && errRes.data.status === 'error'){
                        // message.error(errRes.data.error)
                    }
                });
        // localStorage.removeItem('user');
        // this.props.router.push('/login')
    };
    popoverHide = () => {
        this.setState({
            visible: false,
        });
    };
    handleVisibleChange = (visible) => {
        this.setState({ visible });
    };

    hideUser = (user) => {
        if(!_.isEmpty(user) && user.length === 11){
            return user.substring(0,3) + '****' + user.substring(7,11)
        }
        return user
    }

    render() {
        const { redirectToLogin } = this.state;
        const { responsive, path } = this.props;

        if(redirectToLogin){
            return (
                <Redirect to="/" />
            )
        }

        return (
            <Header style={{ background: '#fff', padding: 0, height: 65 }} className="custom-theme" >
                {
                    responsive.data.isMobile ? (
                        <Popover content={<SiderCustom path={path} popoverHide={this.popoverHide} />} trigger="click" placement="bottomLeft" visible={this.state.visible} onVisibleChange={this.handleVisibleChange}>
                            <Icon type="bars" className="trigger custom-trigger" />
                        </Popover>
                    ) : (
                        <Icon
                            className="trigger custom-trigger"
                            type={this.props.collapsed ? 'menu-unfold' : 'menu-fold'}
                            onClick={this.props.toggle}
                        />
                    )
                }
                <Menu
                    mode="horizontal"
                    style={{ lineHeight: '64px', float: 'right' }}
                    onClick={this.menuClick}
                >
                    <Menu.Item key="full" onClick={this.screenFull} >
                        <Icon type="arrows-alt" onClick={this.screenFull} />
                    </Menu.Item>
                    {/*
                    <Menu.Item key="1">
                        <Badge count={25} overflowCount={10} style={{marginLeft: 10}}>
                            <Icon type="notification" />
                        </Badge>
                    </Menu.Item>
                    */}
                    <SubMenu title={<span className="avatar"><img src={avater} alt="头像" /><i className="on bottom b-white" /></span>}>
                        <MenuItemGroup title="用户中心">
                            <Menu.Item key="setting:1">你好 - {this.hideUser(localStorage.getItem('user'))}</Menu.Item>
                            <Menu.Item key="logout"><span onClick={this.logout}>退出登录</span></Menu.Item>
                        </MenuItemGroup>
                    </SubMenu>
                </Menu>
                <style>{`
                    .ant-menu-submenu-horizontal > .ant-menu {
                        width: 120px;
                        left: -40px;
                    }
                `}</style>
            </Header>
        )
    }
}

const mapStateToProps = state => {
    const { responsive = {data: {}} } = state.httpData;
    return {responsive};
};
const mapDispatchToProps = dispatch => ({
    fetchData: bindActionCreators(fetchData, dispatch),
    receiveData: bindActionCreators(receiveData, dispatch)
});

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(HeaderCustom));
