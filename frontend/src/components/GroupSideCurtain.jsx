import "./SideCurtain.css"
import TextWithInput from "./TextWithInput.jsx";
import React, {useEffect, useMemo, useState} from "react";
import {
    FetchGroupData,
    GroupAddMembers,
    GroupRemMembers,
    SettingGroupBak,
    SettingGroupName,
    SettingNicknameInGroup
} from "../../wailsjs/go/groupbinder/GroupBinder.js";
import {EventsOn} from "../../wailsjs/runtime/runtime.js";
import groupAddIcon from "../assets/images/group_add.png"
import groupSubIcon from "../assets/images/group_substract.png"
import {Button, Checkbox, Input, Modal} from "antd";

import "./GroupSideCurtain.css"

const {TextArea} = Input;


const GroupSideCurtain = ({convItem, messageApi}) => {
    const [groupMebAddModalShow, setGroupMebAddModalShow] = useState(false);
    const [mebAdding, setMebAdding] = useState(false)
    const [addUids, setAddUids] = useState("")

    const [groupMebRemModalShow, setGroupMebRemModalShow] = useState(false);
    const [mebRemoving, setMebRemoving] = useState(false)

    const [selectedIds, setSelectedIds] = useState([])
    const selectedIdSet = useMemo(() => new Set(selectedIds), [selectedIds]);


    useEffect(() => {
        const unSubGroupNicknameModifyEvent = EventsOn("event_backend:modify_group_nickname", groupNo => {
            console.log(`receive [modify_group_nickname], groupNo=${groupNo}, relationId=${convItem.relationId}`);

            if (convItem?.relationId === groupNo) {
                (async () => {
                    await fetchGroupData(groupNo)
                })();
            }
        })

        const unsubGroupMebChangeEvents = EventsOn(
            "event_backend:group_member_changed",
            groupNo => {
                if (convItem?.relationId === groupNo) {
                    (async () => {
                        await fetchGroupData(groupNo)
                    })();
                }
            })

        return () => {
            unSubGroupNicknameModifyEvent()
            unsubGroupMebChangeEvents()
        }
    }, [])

    useEffect(() => {
        // console.log(`useEffect convType=${convItem.convType}`);

        (async () => {
            await fetchGroupData(convItem.relationId)
        })();
        // 拉取群资料
    }, [convItem])

    const [groupData, setGroupData] = useState({})

    const fetchGroupData = async (groupNo) => {
        try {
            let data = await FetchGroupData(groupNo)
            console.log(`${JSON.stringify(data)}`)
            setGroupData(data)
        } catch (e) {
            console.error(e)
        }
    }

    const onGroupNameSet = name => {
        (async () => {
            try {
                await SettingGroupName(convItem.relationId, name)
                // await fetchGroupData()
            } catch (e) {
                console.error(e)
            }
        })()
    }

    const onGroupBakSet = bak => {
        (async () => {
            try {
                await SettingGroupBak(convItem.relationId, bak)
                // await fetchGroupData()
            } catch (e) {
                console.error(e)
            }
        })()
    }

    const onMebItemClick = item => {
        if (item.id > 0) {

        } else {
            if (item.id === -1) {
                // 添加
                setGroupMebAddModalShow(true)
            } else if (item.id === -2) {
                // 移出
                setGroupMebRemModalShow(true)
            }
        }
    }

    const renderGridItem = () => {
        if (!groupData || !groupData.membersInfo) {
            return null
        }

        const displayItems = [...groupData.membersInfo]

        const isMgr = groupData.role < 3;
        if (isMgr) {
            displayItems.push(
                {
                    id: -1,
                    uid: -1,
                    nickname: "添加",
                    role: -1,
                    avatar: groupAddIcon,
                }, {
                    id: -2,
                    uid: -2,
                    nickname: "移出",
                    role: -1,
                    avatar: groupSubIcon,
                }
            );
        }

        return (
            displayItems.map((item, idx) => (
                <div key={item.id ?? idx}
                     onClick={() => onMebItemClick(item)}
                     style={{
                         width: 48,
                         display: "flex",
                         flexDirection: "column",
                         alignItems: "center",
                         overflow: "hidden",
                         boxSizing: "border-box"
                     }}>
                    <img src={item.avatar} style={{
                        width: 34,
                        height: 34,
                        minWidth: 34,
                        minHeight: 34,
                    }}/>

                    <div style={{
                        overflow: "hidden",
                        textOverflow: "ellipsis",
                        whiteSpace: "nowrap",
                        fontSize: 11,
                        width: '100%',
                        color: '#999',
                        marginTop: 3,
                        paddingBottom: 8,
                        boxSizing: "border-box"
                    }}>{item.nickname}
                    </div>
                </div>
            ))
        )
    }


    const onGroupNicknameSet = nickname => {
        (async () => {
            try {
                await SettingNicknameInGroup(convItem.relationId, nickname)
                // await fetchGroupData()
            } catch (e) {
                console.error(e)
            }
        })()
    }


    const renderAddModelFoot = (isAdd) => {
        const groupNo = convItem.relationId

        return (<div style={{
            display: "flex", flexDirection: "row", justifyContent: "end", alignItems: "center", marginTop: 24,
        }}>
            <Button className="chat-modal-btn" style={{marginRight: 16}}
                    onClick={() => {
                        if (isAdd) {
                            setGroupMebAddModalShow(false);
                        } else {
                            setGroupMebRemModalShow(false)
                        }
                    }}>取消</Button>
            <Button
                disabled={isAdd ? mebAdding : mebRemoving}
                type="primary"
                loading={isAdd ? mebAdding : mebRemoving}
                className="chat-modal-btn"
                style={{width: 100}}
                onClick={async () => {
                    if (isAdd) {
                        // 添加群成员
                        setGroupMebAddModalShow(false);

                        setMebAdding(false)

                        try {
                            await GroupAddMembers(groupNo, addUids)
                        } catch (e) {
                            messageApi.error(`group add members failed, groupNo=${groupNo}, uids=${addUids}, e=${e}`)
                        }
                    } else {
                        // 移出群成员
                        setGroupMebRemModalShow(false);

                        setMebRemoving(false);

                        try {
                            await GroupRemMembers(groupNo, selectedIds)
                        } catch (e) {
                            messageApi.error(`group remove members failed, groupNo=${groupNo}, uids=${selectedIds}, e=${e}`)
                        }
                    }
                }}
            >{isAdd ? "添加" : "移出"}</Button>
        </div>)
    }

    const renderAddModalBody = () => {
        return (<div style={{
            boxSizing: "border-box",
            display: "flex",
            flexDirection: "column",
            alignItems: "flex-start"
        }}>
            <TextArea
                value={addUids}
                onChange={e => setAddUids(e.target.value)}
                placeholder="Input Users For Add To Group(Like: u1;u2;u3)"
                style={{
                    boxShadow: 'none',
                    marginTop: 16,
                    height: 120,
                    outline: 'none',
                    caretColor: '#1677FF',
                    resize: 'none'
                }}
            />
        </div>)
    }


    const renderRemModalBody = () => {
        if (!groupData || !groupData.membersInfo) {
            return null
        }

        return (<div style={{
            boxSizing: "border-box",
            display: "flex",
            flexDirection: "column",
            alignItems: "flex-start"
        }}>
            {
                groupData.membersInfo.map((item) => {
                    const isChecked = selectedIdSet.has(item.uid);
                    return (
                        <div
                            key={item.id}
                            className="rem-item"
                            onClick={
                                () => {
                                    const uid = item.uid
                                    if (selectedIdSet.has(uid)) {
                                        setSelectedIds(prev => prev.filter(x => x !== uid));
                                    } else {
                                        setSelectedIds(prev => [...prev, uid]);
                                    }
                                }
                            }
                            style={{
                                minWidth: 316,
                                display: "flex",
                                flexDirection: "row",
                                alignItems: "center",
                                padding: 8,
                            }}>

                            <Checkbox
                                checked={isChecked}
                                value={item.id}/>

                            <img
                                src={item.avatar}
                                style={{
                                    marginLeft: 12,
                                    marginRight: 10,
                                    width: 28,
                                    height: 28,
                                    minWidth: 28,
                                    minHeight: 28,
                                }}/>

                            <span>{item.nickname}</span>
                        </div>
                    )
                })
            }
        </div>)
    }

    return (
        <div>
            <div style={{
                display: "flex",
                flexDirection: "column",
                backgroundColor: "white",
                paddingTop: 18,
            }}>
                <div style={{
                    display: "grid",
                    paddingLeft: 20,
                    overflow: "hidden",
                    paddingRight: 20,
                    gridTemplateColumns: 'repeat(4, 1fr)',
                    gap: 6,
                }}>

                    {renderGridItem()}
                </div>

                <div style={{
                    paddingLeft: 18,
                    paddingRight: 18,
                    height: '100%',
                    display: "flex",
                    flexDirection: "column",
                    alignItems: "start",
                    overflow: "hidden"
                }}>
                    <div style={{
                        height: 1, width: "100%", backgroundColor: '#ddd', marginTop: 10,
                    }}/>

                    <div className="title">群聊名称</div>

                    <TextWithInput
                        text={groupData.groupName}
                        placeholder={"群名称"}
                        onValSetting={onGroupNameSet}
                    />

                    <div className="title">群公告</div>

                    <div
                        className="sub-title">
                        {groupData.groupAnnouncement || "群主暂未设置群公告"}
                    </div>

                    <div className="title">备注</div>

                    <TextWithInput
                        text={groupData.groupBak}
                        placeholder={"群聊的备注仅自己可见"}
                        onValSetting={onGroupBakSet}/>

                    <div className="title">我在本群的昵称</div>

                    <TextWithInput
                        text={groupData.nicknameInGroup}
                        placeholder={"昵称"}
                        onValSetting={onGroupNicknameSet}
                    />
                    <div style={{
                        marginTop: 10,
                        height: 1, width: "100%", backgroundColor: '#ddd'
                    }}/>

                </div>
            </div>

            <Modal
                open={groupMebAddModalShow}
                maskClosable={false}
                title={(<div style={{textAlign: "left"}}>添加群成员 <span
                    style={{color: "#bcbcbc"}}>({groupData.groupBak ?? groupData.groupName})</span></div>)}
                closeIcon={false}
                footer={renderAddModelFoot(true)}
            >
                {renderAddModalBody()}
            </Modal>

            <Modal
                open={groupMebRemModalShow}
                maskClosable={false}
                title={(<div style={{textAlign: "left"}}>移除群成员 <span
                    style={{color: "#bcbcbc"}}>({groupData.groupBak ?? groupData.groupName})</span></div>)}
                closeIcon={false}
                footer={renderAddModelFoot(false)}
            >
                {renderRemModalBody()}
            </Modal>
        </div>
    )
}

export default GroupSideCurtain