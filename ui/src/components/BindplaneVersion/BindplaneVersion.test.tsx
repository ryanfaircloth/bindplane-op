import { render, screen } from "@testing-library/react";
import * as utils from "./utils";
import { BindplaneVersion } from "./BindplaneVersion";

describe("BindplaneVersion", () => {
  beforeEach(() => {
    jest.spyOn(utils, "version").mockReturnValue("1.2.3");
  });
  it("shows return value of utils.version", () => {
    render(<BindplaneVersion />);
    screen.getByText("BindPlane OP 1.2.3");
  });
});
