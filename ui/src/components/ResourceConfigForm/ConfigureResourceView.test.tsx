import { MockedProvider } from "@apollo/client/testing";
import { render, screen } from "@testing-library/react";
import { ConfigureResourceView } from "./ConfigureResourceView";

describe("ConfigureResourceView", () => {
  it("supports pausing destinations", () => {
    const onTogglePause = jest.fn();
    render(
      <MockedProvider>
        <ConfigureResourceView
          kind={"destination"}
          displayName={"Friendly Name"}
          description={"description"}
          paused={false}
          formValues={{}}
          parameterDefinitions={[]}
          onTogglePause={onTogglePause}
        />
      </MockedProvider>
    );
    const togglePause = screen.getByTestId("resource-form-toggle-pause");
    expect(togglePause.textContent).toBe("Pause");
    togglePause.click();
    expect(onTogglePause).toHaveBeenCalledTimes(1);
  });
  it("supports resuming destinations", () => {
    const onTogglePause = jest.fn();
    render(
      <MockedProvider>
        <ConfigureResourceView
          kind={"destination"}
          displayName={"Friendly Name"}
          description={"description"}
          paused={true}
          formValues={{}}
          parameterDefinitions={[]}
          onTogglePause={onTogglePause}
        />
      </MockedProvider>
    );
    const togglePause = screen.getByTestId("resource-form-toggle-pause");
    expect(togglePause.textContent).toBe("Resume");
    togglePause.click();
    expect(onTogglePause).toHaveBeenCalledTimes(1);
  });
});
