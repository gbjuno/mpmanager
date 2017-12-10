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

class PlaceSearch extends Component {
    componentDidMount() {
        // To disabled submit button at the beginning.
        this.props.form.validateFields();
    }
    handleSubmit = (e) => {
        e.preventDefault();
        
        this.props.form.validateFields((err, values) => {
            if (!err) {
                const { fetchData } = this.props
                //fetchData({funcName: 'fetchScPic', stateName: 'picData', params: {picName: values.fileName}});
            }
        });
    };
    render() {
        const { getFieldDecorator, getFieldsError, getFieldError, isFieldTouched } = this.props.form;
        const { style } = this.props

        // Only show error after a field is touched.
        const fileNameError = isFieldTouched('fileName') && getFieldError('fileName');
        return (
            <Form layout="inline" style={style} onSubmit={this.handleSubmit}>
                <FormItem
                    validateStatus={fileNameError ? 'error' : ''}
                    help={fileNameError || ''}
                >
                    {getFieldDecorator('fileName', {
                    })(
                        <Input
                            suffix={<Icon type="search"/>}
                            placeholder="请输入关键字"
                            style={{ width: 200 }}
                            onPressEnter={value => console.log(value)}
                        />
                    )}
                </FormItem>
                <FormItem>
                    <Button
                        type="primary"
                        htmlType="submit"
                    >
                       搜索
                    </Button>
                </FormItem>
            </Form>
        );
    }
}


export default Form.create()(PlaceSearch);