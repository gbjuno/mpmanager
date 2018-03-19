/**
 * Created by Jingle Chen on 2018/3/5.
 */
import React from 'react';
import { Form, Icon, Input, Button, Checkbox, message } from 'antd';
import { connect } from 'react-redux';
import { Link, withRouter, Redirect } from 'react-router-dom';
import { bindActionCreators } from 'redux';
import { fetchData, receiveData } from '@/action';

const FormItem = Form.Item;

class Login extends React.Component {

    state = {
        redirectToHome: false,
    }

    componentWillMount() {
        const { receiveData } = this.props;
        receiveData(null, 'auth');
    }
    componentWillReceiveProps(nextProps) {
        const { auth: nextAuth = {} } = nextProps;
        const { router } = this.props;
        if (nextAuth.data && nextAuth.data.uid) {   // 判断是否登陆
            localStorage.setItem('user', JSON.stringify(nextAuth.data));
            router.push('/');
        }
    }
    handleSubmit = (e) => {
        e.preventDefault();
        this.props.form.validateFields((err, values) => {
            if (!err) {
                console.log('Received values of form: ', {...values});
                const { fetchData } = this.props;
                fetchData({funcName: 'authLogin', params: {...values}, stateName: 'authStatus'})
                .then(res => {
                    if(res.data.status === 200){
                        this.setState({
                            redirectToHome: true,
                        })
                    }
                    message.success('登录成功')
                }).catch(err => {
                    let errRes = err.response
                    if(errRes.data && errRes.data.status === 'error'){
                        message.error(errRes.data.error)
                    }
                });
            }
        });
    };
    gitHub = () => {
        window.location.href = 'https://github.com/login/oauth/authorize?client_id=792cdcd244e98dcd2dee&redirect_uri=http://localhost:3006/&scope=user&state=reactAdmin';
    };
    render() {
        const { getFieldDecorator } = this.props.form;
        const { redirectToHome } = this.state
        const { from } = this.props.location.state || { from: { pathname: "/app/ux/tp" } };

        if(redirectToHome){
            return (
                <Redirect to="/app/ux/tp"/>
            )
        }

        return (
            <div className="login">
                <div className="login-form" >
                    <div className="login-logo">
                        <span>佛山市顺德市安监局</span>
                    </div>
                    <Form onSubmit={this.handleSubmit} style={{maxWidth: '300px'}}>
                        <FormItem>
                            {getFieldDecorator('Phone', {
                                rules: [{ required: true, message: '请输入用户名!' }],
                            })(
                                <Input prefix={<Icon type="user" style={{ fontSize: 13 }} />} placeholder="请输入用户名或手机号码" />
                            )}
                        </FormItem>
                        <FormItem>
                            {getFieldDecorator('Password', {
                                rules: [{ required: true, message: '请输入密码!' }],
                            })(
                                <Input prefix={<Icon type="lock" style={{ fontSize: 13 }} />} type="password" placeholder="请输入密码" />
                            )}
                        </FormItem>
                        <FormItem>
                            <Button type="primary" htmlType="submit" className="login-form-button" style={{width: '100%'}}>
                                登录
                            </Button>
                        </FormItem>
                    </Form>
                </div>
            </div>

        );
    }
}

const mapStateToPorps = state => {
    const { authStatus } = state.httpData;
    return { authStatus };
};
const mapDispatchToProps = dispatch => ({
    fetchData: bindActionCreators(fetchData, dispatch),
    receiveData: bindActionCreators(receiveData, dispatch)
});


export default connect(mapStateToPorps, mapDispatchToProps)(Form.create()(Login));