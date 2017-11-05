/**
 * Created by hao.cheng on 2017/4/15.
 */
import React, { Component } from 'react';

import { Form, Icon, Input, Button } from 'antd';
const FormItem = Form.Item;
const Search = Input.Search;

function hasErrors(fieldsError) {
    return Object.keys(fieldsError).some(field => fieldsError[field]);
}

class SearchForm extends Component {
    componentDidMount() {
        // To disabled submit button at the beginning.
        this.props.form.validateFields();
    }
    handleSubmit = (e) => {
        e.preventDefault();
        this.props.form.validateFields((err, values) => {
            if (!err) {
                console.log('Received values of form: ', values);
            }
        });
    };
    render() {
        const { getFieldDecorator, getFieldsError, getFieldError, isFieldTouched } = this.props.form;
        const { style } = this.props

        // Only show error after a field is touched.
        const userNameError = isFieldTouched('userName') && getFieldError('userName');
        return (
            <Form layout="inline" style={style} onSubmit={this.handleSubmit}>
                <FormItem
                    validateStatus={userNameError ? 'error' : ''}
                    help={userNameError || ''}
                >
                    {getFieldDecorator('userName', {
                        rules: [{ required: true, message: '请输入文件名!' }],
                    })(
                        <Search
                            placeholder="请输入文件名"
                            style={{ width: 200 }}
                            onSearch={value => console.log(value)}
                        />
                    )}
                </FormItem>
                <FormItem>
                    <Button
                        type="primary"
                        htmlType="button"
                    >
                       搜索
                    </Button>
                </FormItem>
            </Form>
        );
    }
}


export default Form.create()(SearchForm);