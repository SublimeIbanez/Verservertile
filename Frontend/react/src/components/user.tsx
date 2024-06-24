import { useState } from "react";
import "./user.css";

export const UserCard = () => {
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

    const handleLogin = (e: React.FormEvent) => {
        e.preventDefault();

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

        if (!hasError) {
            LogIn({ username, password, usernameError, passwordError }, { setPopUp, setUserInfo });
            setPopUp(false);
        }
    };

    return (
        <div className="UserLoginPopup">
            <div className="UserLoginPopup_Inner">
                {
                    usernameError || passwordError ? <p className="InvalidCredentials">Invalid Credentials</p> : <></>
                }
                <div className="UserLoginPopupLabel">
                    <div></div>
                    <h2 className="UserLoginLoginLabel">Login</h2>
                    <button className="UserLoginPopupExit" onClick={() => setPopUp(false)}>X</button>
                </div>
                <form onSubmit={e => handleLogin(e)}>
                    <div className="InputFields">
                        <div className="UsernameLabel">Username:</div>
                        <input className="UsernameInput" type="text" value={username} onChange={e => setUserName(e.target.value)} />
                        <div className="PasswordLabel">Password:</div>
                        <input className="PasswordInput" type="text" value={password} onChange={e => setPassword(e.target.value)} />
                    </div>
                    <button type="submit">Login</button>
                </form>
            </div>
        </div>
    )
}

type LoginProps = {
    username: string;
    password: string;
    usernameError: boolean;
    passwordError: boolean;
}

const LogIn = ({ username, password, usernameError, passwordError }: LoginProps, { setPopUp, setUserInfo }: UserInfoProps) => {
    console.log(username, password);

    setUserInfo(username);
    setPopUp(false);
};

const LogOut = (setUserInfo: React.Dispatch<React.SetStateAction<string>>) => {
    setUserInfo("");
};