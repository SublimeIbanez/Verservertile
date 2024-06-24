import './App.css'
import { UserCard } from './components/user'

function App() {

    return (
        <>
            <div>
                <div className="Header">
                    <div className="Header_Left">
                        This is a logo
                    </div>
                    <div className="Header_Center">
                        IDK have a directory or something here
                    </div>
                    <div className="Header_Right">
                        <UserCard />
                    </div>
                </div>
                <ul style={{
                    listStyleType: "none",
                }}>
                    <li><a href="#">Node Management</a></li>
                    <li><a href="#">Database</a></li>
                </ul>
            </div>
        </>
    )
}

export default App
