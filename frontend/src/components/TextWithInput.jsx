import "./TextWithInput.css"
import {useEffect, useRef, useState} from "react";
import {Input} from "antd";

const TextWithInput = ({text, placeholder, onValSetting}) => {
    const [inputShow, setInputShow] = useState(false)

    const [value, setValue] = useState("")
    useEffect(() => {
        // console.log("text=", text)
        setValue(text || placeholder)
    }, [text])

    const inputRef = useRef(null);

    const onBlur = () => {
        setInputShow(false)

        const val = value
        if (!val) {
            setValue(placeholder)
            return
        }

        onValSetting(val)
        // todo
    }

    const handleKeyDown = (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            inputRef.current?.blur()

        }
    };

    return (<div style={{
        width: '100%', position: "relative", height: 20, display: "flex", alignItems: "center", justifyContent: "center"
    }}>
        <div
            onClick={() => {
                setInputShow(true)
            }}
            className="sub-title">
            {value}
        </div>

        <div style={{
            position: "absolute", left: 0, top: 0, width: '100%', height: 20, display: inputShow ? "flex" : "none"
        }}>
            <Input
                placeholder={placeholder}
                ref={inputRef}
                onKeyDown={handleKeyDown}
                value={value}
                onChange={e => setValue(e.target.value)}
                style={{
                    border: 'none',
                    outline: 'none',
                    background: '#f7f7f7',
                    width: '100%',
                    height: '100%',
                    boxShadow: 'none',
                    padding: 0,
                    color: '#424242',
                    fontSize: 12,
                }}
                autoFocus
                onBlur={onBlur}/>
        </div>
    </div>)
}

export default TextWithInput