import { useState } from "react";

export const User = () => {
    const [userInfo, setUserInfo] = useState("");
    const [popup, setPopUp] = useState(false);

    const togglePopUp = () => {
        setPopUp(!popup);
    };

    return (
        <>
            <div className="UserCard">
                <div className="UserLogIn_Out">
                    {
                        userInfo.length === 0 ?
                            <a onClick={() => togglePopUp()}>Log In</a> :
                            <a onClick={() => LogOut(setUserInfo)}>Log Out</a>
                    }
                    {
                        popup ? <UserInfoPopup setPopUp={setPopUp} setUserInfo={setUserInfo} /> : <></>
                    }
                </div>
                <div className="UserPicture">
                    <img src="vite.svg"></img>
                </div>
            </div>
        </>
    )
};

type UserInfoProps = {
    setPopUp: React.Dispatch<React.SetStateAction<boolean>>;
    setUserInfo: React.Dispatch<React.SetStateAction<string>>;
};

const UserInfoPopup = ({ setPopUp, setUserInfo }: UserInfoProps) => {
    const [username, setUserName] = useState("");
    const [password, setPassword] = useState("");
    const [usernameError, setUsernameError] = useState(false);
    const [passwordError, setPasswordError] = useState(false);

    console.log("kek?");

    const handleLogin = () => {
        let hasError = false;
        if (username.length === 0) {
            setUsernameError(true);
            hasError = true;
        } else {
            setUsernameError(false);
        }
        if (password.length === 0) {
            setPasswordError(true);
            hasError = true;
        } else {
            setPasswordError(false);
        }
        if (hasError) {
            return;
        }
        LogIn(setUserInfo);
        setPopUp(false);
    };

    return (
        <div className="UserLoginPopup">
            <div className="UserLoginPopup_Inner">
                <h2>Login</h2>
                {
                    usernameError ? <p>Bruh... username can't be blank yo</p> : <></>
                }
                {
                    passwordError ? <p>Bruh... password can't be blank yo</p> : <></>
                }
                <form onSubmit={() => handleLogin()}>
                    <label>
                        Username:
                        <input type="text" value={username} onChange={e => setUserName(e.target.value)} />
                    </label>
                    <br></br>
                    <label>
                        Password:
                        <input type="text" value={password} onChange={e => setPassword(e.target.value)} />
                    </label>
                    <br></br>
                    <button type="submit">Login</button>
                </form>
                <button onClick={() => setPopUp(false)}>Close</button>
            </div>
        </div>
    )
}

const LogIn = (setUserInfo: React.Dispatch<React.SetStateAction<string>>) => {
    console.log("You clicked the login kekekekekekekek");

    setUserInfo("User");
};

const LogOut = (setUserInfo: React.Dispatch<React.SetStateAction<string>>) => {
    setUserInfo("");
};