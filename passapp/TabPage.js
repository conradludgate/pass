import React from "react";
import { StyleSheet, Text, View } from "react-native";

import { BarCodeScanner, Permissions } from "expo";

import { Buffer } from "buffer";
import nacl from "tweetnacl";

export default class TabPage extends React.Component {
	state = {
		hasCameraPermission: null,
		value: "",
	};

	async componentWillMount() {
		const { status } = await Permissions.askAsync(Permissions.CAMERA);
		this.setState({ hasCameraPermission: status === "granted" });
	}

	render() {
		if (this.props.page === "1") {
			return <Text>Devices Page</Text>;
		} else if (this.props.page === "2") {
			return <Text>Passwords Page</Text>;
		} else {
			return _pageQR(
				this.state.hasCameraPermission,
				this._handleBarCodeRead.bind(),
			);
		}
	}

	_handleBarCodeRead = ({ type, data }) => {
		bytes = Buffer.from(data, "base64");
		if (bytes.length != 128) return;

		ed = bytes.slice(0, 32);
		curve = bytes.slice(32, 64);

		verify = nacl.sign.detached.verify(
			bytes.slice(0, 64),
			bytes.slice(64),
			ed,
		);

		if (!verify) return;

		// Generate keypair, add it to keepass database and send public keys to server

		this.props.changeTab(1);
	};
}

function _pageQR(perm, callback) {
	if (perm === null) {
		return <Text>Requesting for camera permission</Text>;
	} else if (perm === false) {
		return <Text>No access to camera</Text>;
	} else {
		return (
			<BarCodeScanner
				onBarCodeRead={callback}
				style={StyleSheet.absoluteFill}
			/>
		);
	}
}
