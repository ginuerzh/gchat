import QtQuick 2.0
import Ubuntu.Components 0.1


Column {
    objectName: "loginView"
    id: loginView
    spacing: 10
    property bool logining
    signal login(string username, string password, bool remember)

    Column {
        spacing: 5
        Label {
            text: "Username:"
            fontSize: "large"
        }
        TextField {
            objectName: "usernameInput"
            id: usernameInput
            width: units.gu(33)
            placeholderText: "username@example.com"
            inputMethodHints: Qt.ImhEmailCharactersOnly
            KeyNavigation.tab: passwordInput
            KeyNavigation.down: passwordInput

            onAccepted: {
                if (usernameInput.text.length > 0 && passwordInput.text.length > 0) {
                    loginView.logining = true
                    loginView.login(usernameInput.text, passwordInput.text, rememberPass.checked)
                }
            }
        }
    }

    Column {
        spacing: 5
        Label {
            text: "Password:"
            fontSize: "large"
        }
        TextField {
            objectName: "passwordInput"
            id: passwordInput
            width: units.gu(33)
            echoMode: TextInput.Password
            KeyNavigation.tab: usernameInput

            onAccepted: {
                if (usernameInput.text.length > 0 && passwordInput.text.length > 0) {
                    loginView.logining = true
                    loginView.login(usernameInput.text, passwordInput.text, rememberPass.checked)
                }
            }
        }
    }

    Row {
        spacing: 15
        Row {
            CheckBox {
                id: rememberPass
            }
            Label{
                anchors.verticalCenter: parent.verticalCenter
                text: "Remember Password"
                fontSize: "medium"
            }
        }
        Item {
            width: loginConfirm.width
            height: loginConfirm.height
            Button {
                objectName: "loginConfirm"
                id: loginConfirm
                text: "Login"
                visible: !loginView.logining
                onClicked: {
                    if (usernameInput.text.length > 0 && passwordInput.text.length > 0) {
                        loginView.logining = true
                        loginView.login(usernameInput.text, passwordInput.text, rememberPass.checked)
                    }
                }
            }
            ActivityIndicator{
                anchors.centerIn: parent
                running: loginView.logining
            }
        }
    }
}

