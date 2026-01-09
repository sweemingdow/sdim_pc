import React, {useEffect, useRef, useState} from 'react';
import ChatItem from "./ChatItem.jsx";

const ChatRoom = () => {
    const chatItems = [{
        msgId: 1, msgType: 1, isSelf: true, content: {
            text: "我打打分111我打打分111我打打分111"
        }, sendInfo: {
            sendNickname: "第三个", sendAvatar: "dage.jpg"
        }
    }, {
        msgId: 2, msgType: 1, isSelf: true, content: {
            text: "efds"
        }, sendInfo: {
            sendNickname: "第三个", sendAvatar: "dage.jpg"
        }
    },
        {
            msgId: 3, msgType: 1, isSelf: true, content: {
                text: "我打打分我打打分111我打打分111我打打分111我打打分111我打打分111我打打分111我打打分111我打打分111我打打分111111我打打分111我打打分111"
            }, sendInfo: {
                sendNickname: "第三个", sendAvatar: "dage.jpg"
            }
        },
        {
            msgId: 4, msgType: 1, isSelf: false, content: {
                text: "efds"
            }, sendInfo: {
                sendNickname: "第三个", sendAvatar: "dage.jpg"
            }
        }]
    const [parentSize, setParentSize] = useState({width: 0, height: 0});
    const parentRef = useRef(null);

    useEffect(() => {
        const updateSize = () => {
            if (parentRef.current) {
                const {offsetWidth, offsetHeight} = parentRef.current;

                setParentSize({width: offsetWidth, height: offsetHeight});
            }
        };

        updateSize();

        window.addEventListener('resize', updateSize);

        return () => window.removeEventListener('resize', updateSize);
    }, []);

    return (<div
        ref={parentRef}
        style={{
            flex: 7,
            width: '100%',
            backgroundColor: "#EDEDED", display: "flex", flexDirection: "column"
        }}>

        {chatItems.map(item => (
            <ChatItem key={item.msgId} chat={{item, pw: parentSize.width, ph: parentSize.height}}/>))}


    </div>)
}

export default ChatRoom