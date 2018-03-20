/**
 * Created by Jingle on 2017/12/10.
 */
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import * as _ from 'lodash'
import moment from 'moment';
import { Form, Icon, Input, Button, Select, DatePicker, message } from 'antd';
import { fetchData, receiveData, updateMenu } from '../../action';

const FormItem = Form.Item;
const Search = Input.Search;
const Option = Select.Option;

const dateFormat = 'YYYY-MM-DD';
const queryDateFormat = 'YYYYMMDD';

function hasErrors(fieldsError) {
    return Object.keys(fieldsError).some(field => fieldsError[field]);
}


class MenuForm extends Component {

    state = {
    }

    componentDidMount() {
        // To disabled submit button at the beginning.
        //this.props.form.validateFields();
        const { menu } = this.props
        console.log('control your mind...', menu)
        if(menu !== null && menu !== undefined){
            this.props.form.setFieldsValue({
                name: menu.name,
                url: menu.url,
            })
        }
    }


    handleSubmit = (e) => {
        e.preventDefault();
        const { value } = this.state
        this.props.form.validateFields((err, values) => {
            if (!err) {
                const { fetchData } = this.props
                let saveObj = {
                    name: values.name,
                    company_id: parseInt(values.company_id),
                    monitor_type_id: value.monitor_type_id,
                }
                
            }
        });
    };

    handleChangeName = (e) => {
        const { updateMenu, wechatLocal, menu } = this.props
        console.log('jjjjjj', wechatLocal, e.target.value, menu)
    }

    handleChangeUrl = (e) => {
        if(this.props.onChange){
            this.props.onChange(e.target.value)
        }
    }


    render() {
        const { getFieldDecorator } = this.props.form;
        const { style, filter } = this.props
        const { value } = this.state

        const { menu } = this.props
        let defaultName = '', defaultUrl = ''
        if(menu !== null){
            defaultName = menu.name
            defaultUrl = menu.url
        }

        const formItemLayout = {
            labelCol: {
              xs: { span: 24 },
              sm: { span: 4 },
            },
            wrapperCol: {
              xs: { span: 24 },
              sm: { span: 12 },
            },
          };

        return (
            <Form style={style} onSubmit={this.handleSubmit}>
                <FormItem 
                    {...formItemLayout}
                    style={{}}
                    label="菜单名称"
                >
                    {getFieldDecorator('name', {
                        initialValue: defaultName,
                        rules: [{
                            required: true, message: '请输入菜单名称!',
                        }],
                    })(
                        <Input  onChange={this.handleChangeName}/>
                    )}
                </FormItem>
                <FormItem 
                    {...formItemLayout}
                    style={{}}
                    label="菜单链接"
                >
                    {getFieldDecorator('url', {
                        initialValue: defaultUrl,
                        rules: [{
                            required: true, message: '请输入菜单链接!',
                        }],
                    })(
                        <Input onChange={this.handleChangeUrl}/>
                    )}
                </FormItem>
            </Form>
        );
    }
}

const mapStateToProps = state => {
    const { searchFilter } = state
    return { wechatLocal: state.wechatLocal };
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch),
    updateMenu: bindActionCreators(updateMenu, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(Form.create()(MenuForm))
