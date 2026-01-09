import React, {useState} from 'react';

import {Button, Input} from 'antd';

const {TextArea} = Input;

const ChatInput = ({send}) => {
    const {sendMsg} = send
    const [value, setInputValue] = useState('');

    const onInputValueChanged = (e) => {
        let val = e.target.value

        setInputValue(val)
    }

    const onSendMsg = (_) => {
        const msd = {
            receiver: "test_u2",
            chatType: 1,
            ttl: 0,
            msgContent: {
                type: 1,
                content: {
                    text: value
                }
            }
        }

        sendMsg(msd)
    }

    return (<div style={{
        flex: 2,
        width: '100%',
        backgroundColor: "cyan",
        display: "flex",
        flexDirection: "column",
    }}>

        <div style={{
            height: 50,
            display: "none",
            flexDirection: "column",
            alignItems: "center",
            backgroundColor: "white"
        }}></div>

        <div style={{
            paddingTop: 8,
            paddingLeft: 5,
            paddingRight: 5,
            paddingBottom: 12,
            flex: 1,
            backgroundColor: "#EDEDED",
            display: "flex",
            flexDirection: "column",
        }}>

            <TextArea
                bordered={false}
                style={{
                    flex: 1,
                    border: 0,
                    boxShadow: 'none',
                    outline: 'none',
                    caretColor: '#1677FF',
                    resize: 'none'
                }}
                value={value}
                onChange={onInputValueChanged}
            />
            <Button
                style={{
                    alignSelf: "flex-end",
                    width: 90,
                    height: 32,
                    borderRadius: 4,
                }}
                disabled={!value || value.length <= 0}
                onClick={onSendMsg}
                type="primary">发送(S)</Button>

        </div>


    </div>)
}

export default ChatInput