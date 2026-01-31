import "./SideCurtain.css"
import TextWithInput from "./TextWithInput.jsx";
import {useEffect, useState} from "react";
import {FetchGroupData, SettingGroupName} from "../../wailsjs/go/groupbinder/GroupBinder.js";

const GroupSideCurtain = ({convItem}) => {
    useEffect(() => {
        // console.log(`useEffect convType=${convItem.convType}`);

        (async () => {
            await fetchGroupData()
        })()
        // 拉取群资料
    }, [])

    const [groupData, setGroupData] = useState({})

    const fetchGroupData = async () => {
        try {
            let data = await FetchGroupData(convItem.relationId)
            // console.log(`${JSON.stringify(data)}`)
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

    const onGroupBakSet = name => {
        alert(name)
    }

    const onGroupNicknameSet = name => {
        alert(name)
    }

    return (
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

                {groupData?.membersInfo?.map((item, idx) => (
                    <div key={idx} style={{
                        width: 48,
                        display: "flex",
                        flexDirection: "column",
                        alignItems: "center",
                        overflow: "hidden"
                    }}>
                        <img src={item.avatar} style={{
                            width: 34,
                            height: 34
                        }}/>

                        <div style={{
                            overflow: "hidden",
                            textOverflow: "ellipsis",
                            whiteSpace: "nowrap",
                            fontSize: 11,
                            width: '100%',
                            color: '#999',
                            marginTop: 5,
                            marginBottom: 3,
                        }}>{item.nickname}
                        </div>
                    </div>
                ))}
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
    )
}

export default GroupSideCurtain