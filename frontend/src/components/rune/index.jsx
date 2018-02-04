import React, { Component } from 'react';
import { Popover } from 'antd'


class Rune extends Component{

    state = {
    };

    isUnqualified = () => {
        const { value } = this.props
        if(value.pictures && value.pictures[0]){
            return value.pictures[0].judgement === 'F'
        }
        return false
    }

    getContent = () => {
        const { value } = this.props
        if(value.pictures && value.pictures[0]){
            return value.pictures[0].judgecomment
        }
    }

    render(){
        const isUnqualified = this.isUnqualified()
        const content = this.getContent()
        return (
            <Popover content={content} title="原因" placement="bottom">
                <div style={{display: isUnqualified?'inline':'none'}}className="un-qualified-pic">不合格</div>
            </Popover>
        )
    }
}

export default Rune;
