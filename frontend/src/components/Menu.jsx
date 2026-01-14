const Menu = ({userProfile}) => {
    const {avatar} = userProfile

    return (
        <div style={{
            // backgroundColor: "red",
            display: "flex",
            flexDirection: "column",
            height: '100%',
            boxSizing: "border-box",
        }}>
            <img
                src={avatar}
                id="m-t"
                style={{
                    width: 42,
                    height: 42,
                    flexShrink: 0
                }}/>

            <div id="m-m" style={{
                width: 42,
                flex: 1,
                marginTop: 16,
                // backgroundColor: "cyan",
            }}></div>

            <div id="m-b" style={{
                width: 42,
                height: 100,
                marginTop: 16,
                // backgroundColor: "yellow",
            }}></div>
        </div>
    )
}

export default Menu