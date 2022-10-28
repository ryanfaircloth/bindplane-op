import { BPDestination } from "./destination";

describe("BPResourceConfiguration", () => {
    describe("name", () => {
        it("is the name in destination metadata", () => {
            expect(new BPDestination({
                metadata: { id: "unique-id", name: "name-of-the-destination"},
                spec: {
                    type: "destination",
                    disabled: false,
                },
            }).name()).toBe("name-of-the-destination");
        });
    });

    describe("disabled", () => {
        it("is copied from the destination spec", () => {
            expect(new BPDestination({
                metadata: { id: "id", name: "name"},
                spec: {
                    type: "destination",
                    disabled: true,
                },
            }).spec.disabled).toBe(true);
        });

        it("can be toggled", () => {
            const resource = new BPDestination({
                metadata: { id: "id", name: "name"},
                spec: {
                    type: "destination",
                    disabled: false,
                },
            });
            expect(resource.spec.disabled).toBe(false);

            resource.toggleDisabled()
            expect(resource.spec.disabled).toBe(true);
        });
    });
});
