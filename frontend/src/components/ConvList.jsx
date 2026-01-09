import ConvItem from "./ConvItem.jsx";
import {forwardRef, useEffect, useImperativeHandle, useState} from "react";
import {EventsOn} from "../../wailsjs/runtime/runtime.js";

/*
* 父组件调用子组件的方式, 如果是函数式组件必须通过:
* 在子组件中使用: forwardRef + useImperativeHandle的方式, 暴露出函数给父组件调用
*
* 父组件通过: useRef + ref 来调用
* */
const ConvList = forwardRef(({setMaskShow, setSelectedConvItem}, ref) => {
    const [convItems, setConvItems] = useState([])

    // 暴露给父组件的方法
    useImperativeHandle(ref, () => ({
        // 设置会话列表
        setConvItemsExternal: (items) => {
            setConvItems(items)
        },
    }))

    useEffect(() => {
        const unSubConvListUpdateEvents = EventsOn("event_backend:conv_list_update", items => {
            console.log(`receive conv list update event, items=${JSON.stringify(items)}`)
            setConvItems(items)
        })

        return () => {
            unSubConvListUpdateEvents()
        }
    }, [])

    const handleConvItemClick = (convItem) => {
        // 隐藏遮罩
        if (setMaskShow) {
            setMaskShow(false)
        }
        // 通知父组件选中的会话
        if (setSelectedConvItem) {
            setSelectedConvItem(convItem)
        }
    }

    return (<div>
            {convItems.map((item, idx) => (<ConvItem
                    key={item.convId}
                    idx={idx}
                    convItem={item}
                    onClick={() => handleConvItemClick(item)}
                />))}
        </div>)
})

export default ConvList