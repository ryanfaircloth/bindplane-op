import { MockedProvider } from "@apollo/client/testing";
import { render, screen } from "@testing-library/react";
import { MainView } from "./MainView";

describe("MainView", () => {
    it("supports pausing destinations", () => {
        const onTogglePause = jest.fn();
        render(
            <MockedProvider>
                <MainView
                    kind={"destination"}
                    displayName={"Friendly Name"}
                    description={"description"}
                    paused={false}
                    formValues={{}}
                    parameterDefinitions={[]}
                    onAddProcessor={() => { }}
                    onEditProcessor={() => { }}
                    onRemoveProcessor={() => { }}
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
                <MainView
                    kind={"destination"}
                    displayName={"Friendly Name"}
                    description={"description"}
                    paused={true}
                    formValues={{}}
                    parameterDefinitions={[]}
                    onAddProcessor={() => { }}
                    onEditProcessor={() => { }}
                    onRemoveProcessor={() => { }}
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
