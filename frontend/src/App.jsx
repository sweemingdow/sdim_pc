import './App.css';
import ConvList from "./components/ConvList.jsx";
import Menu from "./components/Menu.jsx";
import SearchAdd from "./components/SearchAdd.jsx";
import TitleSetting from "./components/TitleSetting.jsx";
import ChatRoom from "./components/ChatRoom.jsx";
import ChatInput from "./components/ChatInput.jsx";
import DebugRowOne from "./components/DebugRowOne.jsx";
import {message} from 'antd'
import {Conn2Engine, Disconnect, SendMsg} from "../wailsjs/go/main/App.js";
import {useEffect, useRef, useState} from "react";
import emptyIcon from './assets/images/room_empty.png'
import {SyncConvList} from "../wailsjs/go/syncbinder/SyncBinder.js";
import {UserProfile} from "../wailsjs/go/userbinder/UserBinder.js";
import {FetchNextMsgs} from "../wailsjs/go/msgbinder/MsgBinder.js";
import {EventsOn} from "../wailsjs/runtime/runtime.js";
import {ClearUnreadCount} from "../wailsjs/go/convbinder/ConvBinder.js";
import {StarGroupChat} from "../wailsjs/go/groupbinder/GroupBinder.js";
import SideCurtain from "./components/SideCurtain.jsx";

function App() {
    const [messageApi, contextHolder] = message.useMessage();
    const convListRef = useRef(null);
    /*
    * connState:
    * 1=未连接
    * 2=连接中
    * 3=连接成功
    * 4=连接失败
    * */
    const [connState, setConnState] = useState(1)

    const [maskShow, setMaskShow] = useState(false)

    const [selectedConvItem, setSelectedConvItem] = useState({idx: -1, convItem: null})

    // 未使用ref, 在useEffect中监听事件: event_backend:conv_list_update
    // 有个闭包陷阱
    const selectedConvItemRef = useRef({idx: -1, convItem: null});


    const chatRoomRef = useRef(null);
    const [chatRoomHeight, setChatRoomHeight] = useState(0);

    const chatRoomRootRef = useRef(null);

    const [convItems, setConvItems] = useState([])
    const [currMsgs, setCurrMsgs] = useState([])

    const [userProfile, setUserProfile] = useState({})

    const [sideCurtainShow, setSideCurtainShow] = useState(false)

    useEffect(() => {
        if (!chatRoomRootRef.current) return;

        const updateHeight = () => {
            const container = chatRoomRootRef.current;
            if (container) {
                const rect = container.getBoundingClientRect();
                setChatRoomHeight(rect.height);
            }
        };

        // 初始计算
        updateHeight();

        // 监听窗口变化
        window.addEventListener('resize', updateHeight);

        // 监听 DOM 变化
        const resizeObserver = new ResizeObserver(updateHeight);
        resizeObserver.observe(chatRoomRootRef.current);

        return () => {
            window.removeEventListener('resize', updateHeight);
            if (chatRoomRootRef.current) {
                resizeObserver.unobserve(chatRoomRootRef.current);
            }
        };
    }, []);

    // useEffect(() => {
    //     if (convItems && convItems.length > 0) {
    //         messageApi.info(`会话列表已更新，共 ${convItems.length} 条会话`)
    //     }
    // }, [convItems])

    useEffect(() => {
        const unSubConvListUpdateEvents = EventsOn("event_backend:conv_list_update", data => {
            const items = data.items
            const idx = data.idx

            if (!items) {
                if (idx === -1) {
                    setUserProfile({avatar: ""})
                }
                setConvItems([])
                return;
            }

            // messageApi.info(`receive conv list update event, size=${items.length}, idx=${idx}`)

            const curSelItem = selectedConvItemRef.current;
            const inRoom = curSelItem.idx !== -1 && curSelItem.convItem && curSelItem.convItem.convId === data.items[idx]?.convId

            console.log(`receive conv list update event, size=${items.length}, idx=${idx}, curSelItem.idx=${curSelItem.idx}, inRoom=${inRoom}`)

            setConvItems(items)

            const convItem = data.items[idx]

            if (inRoom) {
                convItem.selected = true
                const newSelItem = {idx: idx, convItem: convItem}
                setSelectedConvItem(newSelItem);
                selectedConvItemRef.current = newSelItem

                setMaskShow(false)

                setCurrMsgs(data.items[idx].recentlyMsgs);
            }
        })

        return () => {
            unSubConvListUpdateEvents()
        }
    }, [])

    useEffect(() => {
        const unSubGlobalHintEvent = EventsOn("event_backend:global_hint", data => {
            // console.log(`receive global_hint event, data=${JSON.stringify(data)}`)

            if (!data) {
                return
            }

            const {errCode, msg, level} = data
            if (level === "info") {
                messageApi.info(`[code=${errCode}]: ${msg}`)
            } else if (level === "warn") {
                messageApi.warning(`[code=${errCode}]: ${msg}`)
            } else if (level === "error") {
                messageApi.error(`[code=${errCode}]: ${msg}`)
            }
        })
        return () => unSubGlobalHintEvent()
    }, [])

    // 主动连接
    const connWithActive = async (uid) => {
        setConnState(2)

        try {
            // 建立连接
            await Conn2Engine(uid)

            setConnState(3)

            messageApi.info(`uid=${uid}, connect success`)

            const convItems = await SyncConvList(uid)

            if (convItems) {
                setConvItems(convItems);

                // if (convListRef.current && convItems) {
                // console.log(`sync conv list success, convItems:${JSON.stringify(convItems)}`);


                // convListRef.current.setConvItemsExternal(convItems);
                // }
            }

            try {
                const up = await UserProfile(uid)
                setUserProfile(up)
            } catch (e) {
                messageApi.error(`fetch profile failed, uid=${uid}, e=${e}`)
            }

            // 会话连接成功, 显示遮罩
            setMaskShow(true)

        } catch (e) {
            setConnState(4)
            messageApi.error(`connect failed, uid=${uid}, e=${e}`)
        }


    }

    // 主动断开连接
    const disConnWithActive = () => {
        setConnState(1)
        setMaskShow(true)

        Disconnect().then(_ => {
            messageApi.warning(`disconnect success`)
        }).catch(e => {
            messageApi.error(`disconnect failed, e=${e}`)
        })
    }


    // 发送消息
    // type MsgSendData struct {
    //     //Sender     string           `json:"sender,omitempty"`   // 发送者uid
    //     Receiver   string           `json:"receiver,omitempty"` // 接收者, 单聊是对方的uid, 群聊是群id
    //     ChatType   preinld.ChatType `json:"chatType,omitempty"`
    //     Ttl        int32            `json:"ttl,omitempty"` // 消息过期时间(sec), -1:阅后即焚,0:不过期
    //     MsgContent *msg.MsgContent  `json:"msgContent,omitempty"`
    // }
    const sendMsg = async (msd) => {
        if (!msd.receiver) {
            if (selectedConvItem && selectedConvItem.idx !== -1 && selectedConvItem.idx < convItems.length) {
                const convItem = convItems[selectedConvItem.idx];
                msd.receiver = convItem.relationId;
                msd.convId = convItem.convId;
                msd.chatType = convItem.convType
            }
        }

        // messageApi.info(`msd=${JSON.stringify(msd)}`);

        try {
            await SendMsg(msd)
            // messageApi.success(`msg send success, msd=${JSON.stringify(msd)}`)

        } catch (e) {
            messageApi.error(`msg send failed, msg=${msd}, e=${e}`)

        }
    }


    const startGroupChat = async (data) => {
        messageApi.info(`start group chat, data=${JSON.stringify(data)}`)

        try {
            await StarGroupChat(data)
        } catch (e) {
            const msg = `group chat create failed, data=${JSON.stringify(data)}, e=${e}`
            console.log(msg)
            messageApi.error(msg)
        }
    }

    const fetchNextMsgs = async () => {
        const convId = selectedConvItem && selectedConvItem.convItem.convId

        console.log(`fetch next msgs, convId=${convId}`)

        try {
            await FetchNextMsgs(convId)
        } catch (e) {
            console.error(`fetch next msgs failed, e=${e}`)
            messageApi.error(`fetch next msgs failed, convId=${convId}, e=${e}`)
        } finally {
            if (chatRoomRef.current) {
                chatRoomRef.current.setMsgLoading(false)
            }
        }
    }

    const onConvItemSelected = (idx, convItem) => {
        // 取消当前选中(反选)
        const cancelCurSelected = selectedConvItem.idx === idx

        // 上次是否选中了条目
        const lastSelected = selectedConvItem.idx !== -1

        if (cancelCurSelected) {// 反选
            const curSelected = !lastSelected
            convItem.selected = curSelected
            const newSelectItem = curSelected ? {idx, convItem} : {idx: -1, convItem: null}

            setSelectedConvItem(newSelectItem);
            selectedConvItemRef.current = newSelectItem

            setMaskShow(!curSelected)

            if (curSelected) {
                // messageApi.info(`on conv item selected, idx:${idx}, convId:${convItem.convId}, relationId=${convItem.relationId}, hasMore=${convItem.hasMore}`)

                if (idx < convItems.length) {
                    // messageApi.info(`msgs=${convItems[idx].recentlyMsgs.length} in convId=${convItem.convId}`)

                    const msgs = convItems[idx].recentlyMsgs

                    setCurrMsgs(msgs);
                }

                // 触发会话清除
                doClearUnread(convItem.convId)
            }

        } else { // 不是反选, 清除之前选择的, 选中当前的
            // selectedConvItem.selected = false
            convItems.forEach(item => {
                item.selected = false
            })

            convItem.selected = true
            const newSelectItem = {idx, convItem}

            setSelectedConvItem(newSelectItem);
            selectedConvItemRef.current = newSelectItem

            setMaskShow(false)

            // messageApi.info(`on conv item newSelected, idx:${idx}, convId:${convItem.convId}, relationId=${convItem.relationId}, hasMore=${convItem.hasMore}`)

            const copyConvItem = {...convItem}
            copyConvItem.recentlyMsgs = []

            // console.log(`selected convItem:\n${JSON.stringify(copyConvItem)}`)

            if (idx < convItems.length) {
                // messageApi.info(`msgs=${convItems[idx].recentlyMsgs.length} in convId=${convItem.convId}`)

                const msgs = convItems[idx].recentlyMsgs

                setCurrMsgs(msgs);
            }

            // 触发会话清除
            doClearUnread(convItem.convId)
        }

        setSideCurtainShow(false)
    }

    const doClearUnread = (convId) => {
        console.log(`trigger conv clear unread in convItem selected, convId=${convId}`);

        (async () => {
            try {
                await ClearUnreadCount(convId)
            } catch (e) {
                messageApi.error(`clear unread failed, convId=${convId}, e=${e}`)
            }
        })()

    }

    const getCurSelectedConvId = () => {
        const curSelItem = selectedConvItemRef.current;
        return curSelItem && curSelItem.idx !== -1 && curSelItem.convItem.convId
    }


    const onTitlePointClick = () => {
        const curSelItem = selectedConvItemRef.current;
        if (curSelItem) {
            // alert(`convType=${curSelItem.convItem.convType}, convId=${curSelItem.convItem.convId}`)
            setSideCurtainShow(!sideCurtainShow)
        }
    }

    return (<>
        {contextHolder}
        <div id="App">
            <div id="outer" style={{
                width: '100vw', height: '100vh', // backgroundColor: 'cyan',
                display: "flex", flexDirection: "column", boxSizing: "border-box", overflow: "hidden"
            }}>
                <div id="dev-top" style={{
                    width: '100%',
                    flex: 1,
                    display: "flex",
                    flexDirection: "column",
                    justifyContent: "flex-start",
                    boxSizing: "border-box"
                }}>
                    <DebugRowOne rowOne={{connState, messageApi, connWithActive, disConnWithActive}}/>

                    <div style={{
                        height: 1, width: "100%", backgroundColor: '#ddd'
                    }}/>

                </div>

                <div id="bottom" style={{
                    flex: 9, display: "flex", flexDirection: "row", boxSizing: "border-box"
                }}>
                    <div id="b-l" style={{
                        backgroundColor: "#eee", height: '100%', boxSizing: "border-box", padding: 10,
                    }}>
                        <Menu userProfile={userProfile}/>
                    </div>

                    <div style={{
                        height: '100%', width: 1, backgroundColor: '#ccc'
                    }}>

                    </div>
                    <div id="b-m" style={{
                        flex: 2.2, minWidth: 0, backgroundColor: '#F7F7F7',
                    }}>
                        <SearchAdd
                            searchAdd={{messageApi, sendMsg, startGroupChat}}/>

                        <ConvList
                            convItems={convItems}
                            ref={convListRef}
                            // setMaskShow={setMaskShow}
                            setSelectedConvItem={onConvItemSelected}/>
                    </div>

                    <div style={{
                        height: '100%', width: 1, backgroundColor: '#ddd'
                    }}>
                    </div>
                    <div id="b-r" style={{
                        flex: 7.3, display: "flex", flexDirection: "column", position: "relative", // backgroundColor:"red",
                    }}>

                        <div style={{
                            boxSizing: "border-box",
                            position: "absolute",
                            width: 248,
                            height: '100%',
                            backgroundColor: "white",
                            right: 0,
                            zIndex: 99,
                            display: 'flex',
                            transform: sideCurtainShow ? 'translateX(0)' : 'translateX(100%)',
                            opacity: sideCurtainShow ? 1 : 0,
                            transition: 'transform 0.3s ease-in-out, opacity 0.3s ease-in-out',
                            pointerEvents: sideCurtainShow ? 'auto' : 'none',
                        }}>
                            <SideCurtain convItem={selectedConvItemRef.current?.convItem} messageApi={messageApi}/>
                        </div>

                        <div style={{
                            position: "absolute",
                            width: '100%',
                            height: '100%',
                            backgroundColor: "#EDEDED",
                            zIndex: 999,
                            display: maskShow ? "flex" : "none",
                            justifyContent: "center",
                            alignItems: "center"
                        }}>
                            <img src={emptyIcon}
                                 style={{
                                     width: 64, height: 64,
                                 }}
                            />
                        </div>

                        <TitleSetting
                            pointClick={onTitlePointClick}
                            setting={{
                                connState,
                                title: selectedConvItem && selectedConvItem.convItem ? selectedConvItem.convItem.title : ""
                            }}/>

                        <div style={{
                            height: 1, width: "100%", backgroundColor: '#ddd'
                        }}/>

                        <div
                            ref={chatRoomRootRef}
                            onClick={() => setSideCurtainShow(false)}
                            style={{
                                flex: 1, backgroundColor: "#EDEDED", overflow: "hidden", // 重要：防止内部滚动干扰
                            }}
                        >

                            <ChatRoom
                                fetchNextMsgs={fetchNextMsgs}
                                height={chatRoomHeight}
                                msgs={currMsgs}
                                ref={chatRoomRef}
                                hasMore={selectedConvItem && selectedConvItem.convItem && selectedConvItem.convItem.hasMore}
                            />
                        </div>


                        <div style={{
                            flexShrink: 0, height: 1, width: "100%", backgroundColor: '#ddd'
                        }}/>

                        <ChatInput
                            send={sendMsg}
                            inputEnabled={connState === 3}
                            onClick={() => {
                                const convId = getCurSelectedConvId()
                                if (convId) {
                                    doClearUnread(convId);
                                }
                            }}
                        />

                    </div>

                </div>

            </div>
        </div>
    </>)
}

export default App