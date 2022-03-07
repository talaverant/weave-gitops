import _ from "lodash";
import React from "react";
import renderer from "react-test-renderer";
import AppContainer from "../App";
import { createMockClient, withContext, withTheme } from "../lib/test-utils";
import { V2Routes } from "../lib/types";

describe("routes", () => {
  _.each(V2Routes, (route) => {
    describe(route, () => {
      it("renders", () => {
        const tree = renderer
          .create(
            withTheme(
              withContext(<AppContainer />, route, {
                applicationsClient: createMockClient({}),
              })
            )
          )
          .toJSON();
        expect(tree).toMatchSnapshot();
      });
    });
  });
});
