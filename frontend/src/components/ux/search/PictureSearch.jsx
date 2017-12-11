/**
 * Created by Jingle on 2017/12/10.
 */
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import * as _ from 'lodash'
import moment from 'moment';
import { Form, Icon, Input, Button, Select, DatePicker } from 'antd';
import { fetchData, receiveData, searchPicture } from '../../../action';

const FormItem = Form.Item;
const Search = Input.Search;
const Option = Select.Option;

const dateFormat = 'YYYY-MM-DD';
const queryDateFormat = 'YYYYMMDD';

function hasErrors(fieldsError) {
    return Object.keys(fieldsError).some(field => fieldsError[field]);
}

class PictureSearch extends Component {

    state = {
        townsData: [],
        villagesData: [],
        selectedTownId: '',
        selectedDate: new Date(),
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
                //fetchData({funcName: 'fetchScPic', stateName: 'picData', params: {picName: values.fileName}});
                console.log('cccsssss', values)
                searchPicture({date: values.selectedDate.format(queryDateFormat)});
                fetchData({funcName: 'fetchPicturesWithPlace', params: { day: values.selectedDate.format(queryDateFormat)}, 
                    stateName: 'picturesData'})
            }
        });
    };

    onDateChange = (date, dateString) => {
        const { searchPicture } = this.props
        if (date === undefined || date === null) return
        
    }

    onTownChange = (value) => {
        this.setState({
            selectedTownId: value,
        })

    }

    

    fetchTownList = () => {
        const { fetchData } = this.props
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

    fetchVillageList = () => {
        const { fetchData } = this.props
        fetchData({funcName: 'fetchCountries', stateName: 'villagesData'}).then(res => {
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
            return <Option key={item.key} value={`${item.id}`}>{item.name}</Option>
        })
    }

    render() {
        const { getFieldDecorator, getFieldsError, getFieldError, isFieldTouched } = this.props.form;
        const { style, filter } = this.props
        const { townsData, selectedDate } = this.state

        console.log('filter .... search cccc', filter)

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
                        placeholder="请选择村"
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
                        placeholder="请选择公司"
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
                        initialValue: moment(selectedDate, dateFormat)
                    })(
                        <DatePicker onChange={this.onDateChange}/>
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
    const { searchFilter } = state
    console.log('pic state------======>>>>>', state)
    return { ...state.httpData, filter: searchFilter};
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch),
    searchPicture: bindActionCreators(searchPicture, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(Form.create()(PictureSearch))
