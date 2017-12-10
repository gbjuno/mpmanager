/**
 * Created by Jingle on 2017/12/10.
 */
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import * as _ from 'lodash'
import { fetchData, receiveData } from '../../../action';
import { Form, Icon, Input, Button, Select, DatePicker } from 'antd';
const FormItem = Form.Item;
const Search = Input.Search;
const Option = Select.Option;

function hasErrors(fieldsError) {
    return Object.keys(fieldsError).some(field => fieldsError[field]);
}

class PictureSearch extends Component {

    state = {
        townsData: [],
        villagesData: [],
        selectedTown: '',
    }

    componentDidMount() {
        // To disabled submit button at the beginning.
        this.props.form.validateFields();
        this.fetchTownList();
    }

    

    handleSubmit = (e) => {
        e.preventDefault();
        
        this.props.form.validateFields((err, values) => {
            if (!err) {
                const { fetchData } = this.props
                fetchData({funcName: 'fetchScPic', stateName: 'picData', params: {picName: values.fileName}});
            }
        });
    };

    onDateChange = () => {

    }

    onTownChange = (value) => {
        this.setState({
            selectedTown: value,
        })
    }

    fetchTownList = () => {
        const { fetchData } = this.props
        console.log('search picture', this.props)
        fetchData({funcName: 'fetchTowns', stateName: 'townsData'}).then(res => {
            if(res === undefined || res.data === undefined || res.data.towns === undefined) return
            this.setState({
                townsData: [...res.data.towns.map(val => {
                    val.key = val.id;
                    return val;
                })],
                loading: false,
            });
        });
    }

    fetchCountriesData = () => {
        if (this.state.selectedTown === undefined) return
        const { fetchData } = this.props
        fetchData({funcName: 'fetchCountries', stateName: 'villagesData', 
            params: {townId: this.state.selectedTown}}).then(res => {
            this.setState({
                villagesData: [...res.data.countries.map(val => {
                    val.key = val.id;
                    return val;
                })],
                loading: false,
            });
        }).catch(err => {
            this.setState({
                villagesData: [],
            })
        });
    }

    getTownOptions = ( townsData=[] ) => {
        
        return townsData.map(item => {
            return <Option key={item.key} value={item.id}>{item.name}</Option>
        })
    }

    render() {
        const { getFieldDecorator, getFieldsError, getFieldError, isFieldTouched } = this.props.form;
        const { style } = this.props
        const { townsData } = this.state

        // Only show error after a field is touched.
        const fileNameError = isFieldTouched('fileName') && getFieldError('fileName');
        return (
            <Form layout="inline" style={style} onSubmit={this.handleSubmit}>
                <FormItem
                    validateStatus={fileNameError ? 'error' : ''}
                    help={fileNameError || ''}
                >
                    {getFieldDecorator('town', {
                    })(
                        <Select
                        showSearch
                        style={{ width: 200 }}
                        placeholder="请选择镇"
                        optionFilterProp="children"
                        onChange={this.onTownChange}
                        filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                        >
                            {this.getTownOptions(townsData)}
                        </Select>
                    )}
                </FormItem>
                <FormItem
                    validateStatus={fileNameError ? 'error' : ''}
                    help={fileNameError || ''}
                >
                    {getFieldDecorator('village', {
                    })(
                        <Select
                        showSearch
                        style={{ width: 200 }}
                        placeholder="请选择镇"
                        optionFilterProp="children"
                        onChange={this.onTownChange}
                        filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                        >
                            {this.getTownOptions(townsData)}
                        </Select>
                    )}
                </FormItem>
                <FormItem
                    validateStatus={fileNameError ? 'error' : ''}
                    help={fileNameError || ''}
                >
                    {getFieldDecorator('company', {
                    })(
                        <Select
                        showSearch
                        style={{ width: 200 }}
                        placeholder="请选择镇"
                        optionFilterProp="children"
                        onChange={this.onTownChange}
                        filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                        >
                            {this.getTownOptions(townsData)}
                        </Select>
                    )}
                </FormItem>
                <FormItem
                    validateStatus={fileNameError ? 'error' : ''}
                    help={fileNameError || ''}
                >
                    {getFieldDecorator('selectedDate', {
                    })(
                        <DatePicker onChange={this.onDateChange} />
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

const mapStateToProps = state => {
    return { ...state.httpData };
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch)
});

export default connect(mapStateToProps, mapDispatchToProps)(Form.create()(PictureSearch))
