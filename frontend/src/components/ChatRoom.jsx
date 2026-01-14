import React, {forwardRef, useEffect, useImperativeHandle, useRef, useState} from 'react';
import ChatItem from "./ChatItem.jsx";
import InfiniteScroll from 'react-infinite-scroll-component';
import {Spin} from "antd";
import './ChatRoom.css'
import {prettyTime} from "../utils/time_format.js";

/*const testChatItems = [{
    msgId: 1, msgType: 1, isSelf: true, content: {
        text: "我打打分111我打打分111我打打分111"
    }, sendInfo: {
        sendNickname: "第三个", sendAvatar: "dage.jpg"
    }
}, {
    msgId: 2, msgType: 1, isSelf: true, content: {
        text: "我打打分111我打打分111我打打分111"
    }, sendInfo: {
        sendNickname: "第三个", sendAvatar: "dage.jpg"
    }
}, {
    msgId: 3, msgType: 1, isSelf: true, content: {
        text: "我打打分111我打打分111我打打分11155"
    }, sendInfo: {
        sendNickname: "第三个", sendAvatar: "dage.jpg"
    }
}, {
    msgId: 4, msgType: 1, isSelf: true, content: {
        text: "我打打分111我打打分111我打打分111"
    }, sendInfo: {
        sendNickname: "第三个", sendAvatar: "dage.jpg"
    }
}, {
    msgId: 5, msgType: 1, isSelf: true, content: {
        text: "我打打分111我打打分111我打打分111"
    }, sendInfo: {
        sendNickname: "第三个", sendAvatar: "dage.jpg"
    }
}, {
    msgId: 6, msgType: 1, isSelf: true, content: {
        text: "efds"
    }, sendInfo: {
        sendNickname: "第三个", sendAvatar: "dage.jpg"
    }, timeShow: true,
}, {
    msgId: 7, msgType: 1, isSelf: true, content: {
        text: "我打打分我打打分111我打打分111我打打分111我打打分111我打打分111我打打分111我打打分111我打打分111我打打分111111我打打分111我打打分111"
    }, sendInfo: {
        sendNickname: "第三个", sendAvatar: "dage.jpg"
    }
}, {
    msgId: 8, msgType: 1, isSelf: false, content: {
        text: "我打打分我打打分111我打打分111我打打分111我打打分111我打打分111我打打分111我打打分111我打打分111我打打分111111我打打分111我打打分111efds我打打分我打打分111我打打分111我打打分111我打打分111我打打分111我打打分111我打打分111我打打分111我打打分111111我打打分111我打打分111"
    }, sendInfo: {
        sendNickname: "第三个", sendAvatar: "dage.jpg"
    }
}]*/

const ChatRoom = forwardRef(({height, msgs, hasMore, fetchNextMsgs}, ref) => {
    // const [chatItems, setChatItems] = useState([])
    const [loading, setLoading] = useState(false)
    console.log(`chat room hasMore=${hasMore}`)

    const [parentSize, setParentSize] = useState({width: 0, height: 0});
    const parentRef = useRef(null);

    useEffect(() => {
        const updateSize = () => {
            if (parentRef.current) {
                const {offsetWidth, offsetHeight} = parentRef.current;

                // console.log("size changed", offsetWidth, offsetHeight)

                setParentSize({width: offsetWidth, height: offsetHeight});
            }
        };

        updateSize();

        window.addEventListener('resize', updateSize);

        return () => window.removeEventListener('resize', updateSize);
    }, []);


    useImperativeHandle(ref, () => (
        {
            setMsgLoading: () => {
                setLoading(false)
            }
        }
    ))


    useEffect(() => {
        if (parentRef.current) {
            // 因为用了 flex-direction: column-reverse，
            // 最新消息在顶部，所以滚动到顶部即“底部”
            // parentRef.current.scrollTo(0, 0);

            const {scrollTop} = parentRef.current;
            console.log(`scrollTop=${scrollTop}`)

            if (scrollTop >= -100) {
                parentRef.current.scrollTo(0, 0);
            }
        }
    }, [msgs]);

    const fetchHistoryMsg = () => {
        if (loading || !hasMore) return;

        setLoading(true);
        console.log("fetchHistoryMsg");
        fetchNextMsgs()

        // setTimeout(() => {
        //     const lastMsgId = chatItems[chatItems.length - 1].msgId
        //     const oneMsgId = lastMsgId + 1
        //     const twoMsgId = oneMsgId + 1
        //     const newHistoryMessages = [{
        //         msgId: oneMsgId, msgType: 1, isSelf: true, content: {
        //             text: "这是历史消息1，加载到顶部"
        //         }, sendInfo: {
        //             sendNickname: "历史用户1", sendAvatar: "dage.jpg"
        //         }
        //     }, {
        //         msgId: twoMsgId, msgType: 1, isSelf: false, content: {
        //             text: "这是历史消息2，加载到顶部"
        //         }, sendInfo: {
        //             sendNickname: "历史用户2", sendAvatar: "dage.jpg"
        //         }
        //     }];
        //
        //     // 将历史消息添加到列表开头
        //     setChatItems(prev => [...prev, ...newHistoryMessages,]);
        //
        //     setLoading(false);
        //     setHasMore(true)
        //
        // }, 2000);
    };

    const now = Date.now()
    const TIME_THRESHOLD = 5 * 60 * 1000;
    const SHORT_GAP = 30 * 1000;

    const processMessages = messages => {
        if (!messages || messages.length === 0) {
            return [];
        }

        // 处理时间显示
        return messages.map((msg, index) => {
            // 第一条消息总是显示时间
            if (index === 0) {
                return {
                    ...msg,
                    timeText: prettyTime(msg.cts),
                    timeShow: true,
                };
            }

            const prevMsg = messages[index - 1];

            // 时间间隔
            const timeDiff = Math.abs(msg.cts - prevMsg.cts);

            // 如果与前一条消息间隔超过5分钟，显示时间
            if (timeDiff > TIME_THRESHOLD) {
                return {
                    ...msg,
                    timeText: prettyTime(msg.cts),
                    timeShow: true,
                };
            }

            // 如果前一条消息显示了时间，且间隔小于30秒，当前消息不显示时间
            if (prevMsg.timeShow && timeDiff <= SHORT_GAP) {
                return {
                    ...msg,
                    timeShow: false,
                };
            }

            // 前一条消息没显示时间，但间隔大于30秒，显示时间
            if (!prevMsg.timeShow && timeDiff > SHORT_GAP) {
                return {
                    ...msg,
                    timeText: prettyTime(msg.cts),
                    timeShow: true,
                };
            }

            // 其他情况不显示时间
            return {
                ...msg,
                timeShow: false,
            };
        })
    };


    return (<div
        id="roomRoot"
        ref={parentRef}
        style={{
            width: '100%', height: height,
            // height:'500',
            overflow: "auto",
            // backgroundColor:"#EDEDED",
            display: "flex",
            flexDirection: "column-reverse",
        }}>

        <InfiniteScroll
            scrollableTarget="roomRoot"
            next={fetchHistoryMsg}
            pullDownToRefresh={false}
            inverse={true}
            hasMore={hasMore}
            loader={
                <div style={{paddingTop: 10, paddingBottom: 10}}>
                    <Spin size="small"/>
                </div>
            }
            dataLength={msgs.length}
            style={{
                display: 'flex', flexDirection: 'column-reverse', // 配合 inverse 使用
                minHeight: '100%',
            }}
        >

            {

                processMessages(msgs).map(msg => (
                    <ChatItem key={msg.msgId || msg.clientId}
                              chat={{msg, pw: parentSize.width, ph: parentSize.height}}/>))
            }

        </InfiniteScroll>


    </div>)
})

export default ChatRoom