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


    const handleKeyDown = (e) => {
        if (e.key === 'Enter') {
            if (e.shiftKey) {
                // Shift + Enter：允许默认行为（换行）
                return;
            }

            // 仅 Enter：阻止换行，触发发送
            e.preventDefault();
            onSendMsg();
        }
    };

    const onSendMsg = (_) => {
        const msd = {
            chatType: 1,
            ttl: 0,
            msgContent: {
                type: 1,
                content: {
                    text: value
                }
            }
        }

        setInputValue("")
        sendMsg(msd)
    }

    return (<div style={{
        width: '100%',
        height: 180,
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
                onKeyDown={handleKeyDown}
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