/**
 * Created by hao.cheng on 2017/4/13.
 */
import React, { Component } from 'react';
import { Layout, Menu, Icon } from 'antd';
import { Link, withRouter } from 'react-router-dom';
const { Sider } = Layout;
const SubMenu = Menu.SubMenu;

class SiderCustom extends Component {
    state = {
        collapsed: false,
        mode: 'inline',
        openKey: '',
        selectedKey: '',
        firstHide: true,        // 点击收缩菜单，第一次隐藏展开子菜单，openMenu时恢复
    };
    componentDidMount() {
        this.setMenuOpen(this.props);
    }
    componentWillReceiveProps(nextProps) {
        this.onCollapse(nextProps.collapsed);
        this.setMenuOpen(nextProps)
    }
    setMenuOpen = props => {
        const {pathname} = props.location;
        this.setState({
            openKey: pathname.substr(0, pathname.lastIndexOf('/')),
            selectedKey: pathname
        });
    };
    onCollapse = (collapsed) => {
        this.setState({
            collapsed,
            firstHide: collapsed,
            mode: collapsed ? 'vertical' : 'inline',
        });
    };
    menuClick = e => {
        this.setState({
            selectedKey: e.key
        });
        const { popoverHide } = this.props;     // 响应式布局控制小屏幕点击菜单时隐藏菜单操作
        popoverHide && popoverHide();
    };
    openMenu = v => {
        this.setState({
            openKey: v[v.length - 1],
            firstHide: false,
        })
    };
    render() {
        return (
            <Sider
                trigger={null}
                breakpoint="lg"
                collapsed={this.props.collapsed}
                style={{overflowY: 'auto'}}
            >
                <div className="logo" />
                <Menu
                    onClick={this.menuClick}
                    theme="dark"
                    mode="inline"
                    selectedKeys={[this.state.selectedKey]}
                    openKeys={this.state.firstHide ? null : [this.state.openKey]}
                    onOpenChange={this.openMenu}
                >
                    {/*
                    <Menu.Item key="/app/dashboard/index">
                        <Link to={'/app/dashboard/index'}><Icon type="mobile" /><span className="nav-text">首页</span></Link>
                    </Menu.Item>
                    */}
                    <Menu.Item key="/app/ux/tp"><Link to={'/app/ux/tp'}><Icon type="picture" /><span className="nav-text">图片管理</span></Link></Menu.Item>
                    <Menu.Item key="/app/ux/cz"><Link to={'/app/ux/cz'}><Icon type="appstore-o" /><span className="nav-text">村镇管理</span></Link></Menu.Item>
                    <Menu.Item key="/app/ux/gs"><Link to={'/app/ux/gs'}><Icon type="home" /><span className="nav-text">公司管理</span></Link></Menu.Item>
                    <Menu.Item key="/app/ux/yh"><Link to={'/app/ux/yh'}><Icon type="user" /><span className="nav-text">用户管理</span></Link></Menu.Item>
                    <Menu.Item key="/app/ux/dd"><Link to={'/app/ux/dd'}><Icon type="environment-o" /><span className="nav-text">地点管理</span></Link></Menu.Item>
                    <Menu.Item key="/app/ux/jq"><Link to={'/app/ux/jq'}><Icon type="environment-o" /><span className="nav-text">假期管理</span></Link></Menu.Item>
                    <Menu.Item key="/app/ux/wwc"><Link to={'/app/ux/wwc'}><Icon type="line-chart" /><span className="nav-text">完成率统计</span></Link></Menu.Item>
                    <Menu.Item key="/app/ux/tj"><Link to={'/app/ux/tj'}><Icon type="line-chart" /><span className="nav-text">统计报表</span></Link></Menu.Item>
                    {/* <SubMenu
                        key="/app/wechat"
                        title={<span><Icon type="wechat" /><span className="nav-text">微信管理</span></span>}
                    >
                        <Menu.Item key="/app/wechat/cd"><Link to={'/app/wechat/cd'}><Icon type="bars" /><span className="nav-text">菜单管理</span></Link></Menu.Item>
                        <Menu.Item key="/app/wechat/wz"><Link to={'/app/wechat/wz'}><Icon type="file-text" /><span className="nav-text">文章管理</span></Link></Menu.Item>
                        <Menu.Item key="/app/wechat/sc"><Link to={'/app/wechat/sc'}><Icon type="cloud-o" /><span className="nav-text">素材管理</span></Link></Menu.Item>
                        <Menu.Item key="/app/wechat/xx"><Link to={'/app/wechat/xx'}><Icon type="message" /><span className="nav-text">消息管理</span></Link></Menu.Item>
                        <Menu.Item key="/app/wechat/ym"><Link to={'/app/wechat/ym'}><Icon type="profile" /><span className="nav-text">页面模板</span></Link></Menu.Item>
                    </SubMenu> */}
                </Menu>
                <style>
                    {`
                    #nprogress .spinner{
                        left: ${this.state.collapsed ? '70px' : '206px'};
                        right: 0 !important;
                    }
                    `}
                </style>
            </Sider>
        )
    }
}

export default withRouter(SiderCustom);