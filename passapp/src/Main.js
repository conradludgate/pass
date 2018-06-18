import React from "react";
import { StyleSheet, Text, View } from "react-native";

import Tabs from "react-native-tabs";
import TabPage from "./TabPage";

export default class Main extends React.Component {
	state = { 
		page: "1",
	};

	constructor(props) {
		super(props);
	}

	render() {
		return (
			<View style={styles.container}>
				<Tabs
					selected={this.state.page}
					onSelect={this._onTabSelect}
					style={{ backgroundColor: "white" }}
					selectedStyle={{ color: "red" }} >

					<Text name="0">Scan</Text>
					<Text name="1">Devices</Text>
					<Text name="2">Passwords</Text>
				</Tabs>

				<TabPage page={this.state.page} changeTab={this._changeTab} />
			</View>
		);
	}

	_changeTab = tab => {
		tabs = ["0", "1", "2"];
		if (tab > 0 && tab < tabs.length) {
			this.setState({ page: tabs[tab] });
		}
	};

	_onTabSelect = tab => {
		this.setState({ page: tab.props.name });
	};
}

styles = StyleSheet.create({
	container: {
		flex: 1,
		justifyContent: "center",
		alignItems: "center",
		backgroundColor: "#F5FCFF",
	},
	welcome: {
		fontSize: 20,
		textAlign: "center",
		margin: 10,
	},
	instructions: {
		textAlign: "center",
		color: "#333333",
		marginBottom: 5,
	},
});
