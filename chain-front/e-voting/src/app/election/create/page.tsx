"use client"

import "primereact/resources/themes/lara-light-cyan/theme.css";
import {Button} from "primereact/button";
import {InputText} from "primereact/inputtext";
import * as yup from "yup";
import {useFormik} from "formik";
import {Message} from "primereact/message";
import {Dropdown} from "primereact/dropdown";
import {Calendar} from "primereact/calendar";
import {Stepper} from "primereact/stepper";
import {useRef} from "react";
import {StepperPanel} from "primereact/stepperpanel";
import {InputNumber} from "primereact/inputnumber";
import {useMutation, useQueryClient} from "@tanstack/react-query";
import axios from "axios";
import {useRouter} from "next/navigation";

interface ElectionData {
    authMethod: string;
    description: string;
    endDate: string;
    numberOfAuthAttempts: string;
    numberOfCandidates: string;
    numberOfSelection: string;
    startDate: string;
    title: string;
    type: string;
}

interface ElectionResponse {
    response: string;
}

const createElection = async (data: ElectionData): Promise<ElectionResponse> => {
    const { data: response } = await axios.post('http://localhost:3000/elections/create', data);
    return response.data;
}
export default function CreateElection() {
    const queryClient = useQueryClient();
    const router = useRouter();

    const { mutate, isLoading } = useMutation({
        mutationFn: createElection,
        onSuccess: () => {
            router.push('/election');
        }
    })

    const validationSchema = yup.object({
        title: yup.string().required('Title is required'),
        description: yup.string().required('Description is required'),
        type: yup.object().nullable().required('Type is required'),
        authMethod: yup.object().nullable().required('Authentication method is required'),
        numberOfAuthAttempts: yup.string().required('Number of authentication attempts is required'),
        numberOfSelection: yup.string().required('Maximum number of selections is required'),
        numberOfCandidates: yup.string().required('The number of candidates'),
        startDate: yup.string().required("Start date is required"),
        endDate: yup.string()
            .test(
                "dates-test",
                "End date should be after start date",
                (value, context) => {
                    let startDate = new Date(context.parent.startDate);
                    let endDate = new Date(context.parent.endDate);
                    return endDate.getTime() > startDate.getTime()
                }
            )
            .required("End date is required")
    });

    const formik= useFormik({
        initialValues: {
            title: '',
            description: '',
            type: null,
            startDate: null,
            endDate: null,
            authMethod: null,
            numberOfAuthAttempts: null,
            numberOfSelection: null,
            numberOfCandidates: null
        },

        validationSchema,

        onSubmit: (values) => {
            const formattedValues = {
                ...values,
                numberOfAuthAttempts: values.numberOfAuthAttempts.toLocaleString(),
                numberOfSelection: values.numberOfSelection.toLocaleString(),
                numberOfCandidates: values.numberOfCandidates.toLocaleString(),
                startDate: values.startDate.toLocaleString(),
                endDate: values.endDate.toLocaleString(),
                authMethod: values.authMethod.name.toLocaleString(),
                type: values.type?.name // Extract the `code` property
            };
            mutate(formattedValues)
        }
    })

    const electionTypes = [
        {name: 'General elections', code: 'GEL'},
        {name: 'Referendum', code: 'RFD'},
        {name: 'Internal governance vote', code: 'VCI'}
    ]

    const authenticationMethods = [
        {name: 'Unique id'},
        {name: 'Email + password'},
        {name: 'OTP'}
    ]

    const stepperRef = useRef(null);

    return (
            <div className="surface-0 p-2">
                <div className="flex justify-content-start">
                    <Button
                        severity="secondary"
                        icon="pi pi-angle-left"
                        label="Elections"
                        raised
                        onClick={() => router.push('/election')}
                    />
                </div>
                <div className="flex align-items-center justify-content-center">
                    <div className="surface-card p-4 shadow-2 border-round w-full lg:w-6">
                        <div className="text-center mb-5">
                            <div className="text-900 text-3xl font-medium mb-3">Election setup form</div>
                        </div>
                        <div>
                            <form onSubmit={formik.handleSubmit}>
                                <Stepper ref={stepperRef} style={{flexBasis: '50rem'}} orientation="vertical">
                                    <StepperPanel header="General information about the elections">
                                        <div className="flex flex-column">
                                            <div className="border-2 border-dashed surface-border border-round surface-ground flex-auto flex justify-content-center align-items-center font-medium p-3">
                                                <div>
                                                    <label htmlFor="title" className="block text-900 font-medium mb-2 mt-2">Election
                                                        title</label>
                                                    <InputText id="title" type="text" placeholder="Election title"
                                                               className="w-full mb-2" name="title"
                                                               onChange={formik.handleChange}
                                                               value={formik.values.title}
                                                    />
                                                    {formik.errors.title && formik.touched.title &&
                                                        <Message severity={"error"} text={formik.errors.title}
                                                                 className="w-full mb-1"/>}

                                                    <label htmlFor="description"
                                                           className="block text-900 font-medium mb-2">Description</label>
                                                    <InputText id="description" type="text" placeholder="Description"
                                                               className="w-full mb-3" name="description"
                                                               value={formik.values.description}
                                                               onChange={formik.handleChange}/>
                                                    {formik.errors.description && formik.touched.description &&
                                                        <Message severity={"error"} text={formik.errors.description}
                                                                 className="w-full mb-1"/>}

                                                    <Dropdown name="type" options={electionTypes} optionLabel="name"
                                                              placeholder="Select a type" className="w-full mb-2"
                                                              value={formik.values.type}
                                                              onChange={(e) => formik.setFieldValue('type', e.value)}
                                                    />
                                                    {formik.errors.type && formik.touched.type &&
                                                        <Message severity={"error"} text={formik.errors.type}
                                                                 className="w-full mb-1"/>}

                                                    <div className="formgrid grid">
                                                        <div className="field col">
                                                            <label htmlFor="startDate"
                                                                   className="block text-900 font-medium mb-2">Start
                                                                date</label>
                                                            <Calendar name="startDate" value={formik.values.startDate}
                                                                      className="w-full"
                                                                      onChange={formik.handleChange}/>
                                                            {formik.errors.startDate && formik.touched.startDate &&
                                                                <Message severity={"error"} text={formik.errors.startDate}
                                                                         className="w-full mt-1"/>}
                                                        </div>
                                                        <div className="field col">
                                                            <label htmlFor="endDate"
                                                                   className="block text-900 font-medium mb-2">End
                                                                date</label>
                                                            <Calendar name="endDate" value={formik.values.endDate}
                                                                      className="w-full"
                                                                      onChange={formik.handleChange}/>
                                                            {formik.errors.endDate && formik.touched.endDate &&
                                                                <Message severity={"error"} text={formik.errors.endDate}
                                                                         className="w-full mt-1"/>}
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                            <div className="flex pt-4 justify-content-end">
                                                <Button label="Next" icon="pi pi-arrow-right" iconPos="right" type="button" raised
                                                        onClick={() => stepperRef.current.nextCallback()}/>
                                            </div>
                                        </div>
                                    </StepperPanel>
                                    <StepperPanel header="Authentication settings">
                                        <div className="flex flex-column">
                                            <div className="border-2 border-dashed surface-border border-round surface-ground flex-auto flex justify-content-center align-items-center font-medium p-3">
                                                <div className="w-full">
                                                    <Dropdown name="authMethod" options={authenticationMethods} optionLabel="name"
                                                              placeholder="Select Auth Method" className="w-full mb-2"
                                                              value={formik.values.authMethod}
                                                              onChange={(e) => formik.setFieldValue('authMethod', e.value)}
                                                    />
                                                    {formik.errors.authMethod && formik.touched.authMethod &&
                                                        <Message severity={"error"} text={formik.errors.authMethod}
                                                                 className="w-full mb-1"/>}

                                                    <InputNumber id="numberOfAuthAttempts"  placeholder="Number of authentication attempts"
                                                               className="w-full mb-3" name="numberOfAuthAttempts"
                                                               value={formik.values.numberOfAuthAttempts}
                                                               onValueChange={formik.handleChange} showButtons min={0} max={100} mode="decimal" step={1}
                                                    />
                                                    {formik.errors.numberOfAuthAttempts && formik.touched.numberOfAuthAttempts &&
                                                        <Message severity={"error"} text={formik.errors.numberOfAuthAttempts}
                                                                 className="w-full mb-1"/>}
                                                </div>
                                            </div>
                                        </div>
                                        <div className="flex pt-4 justify-content-between">
                                            <Button label="Back" severity="secondary" icon="pi pi-arrow-left" raised
                                                    onClick={() => stepperRef.current.prevCallback()} type="button"/>
                                            <Button label="Next" icon="pi pi-arrow-right" iconPos="right" raised
                                                    onClick={() => stepperRef.current.nextCallback()} type="button"/>
                                        </div>
                                    </StepperPanel>
                                    <StepperPanel header="Defining voting options">
                                        <div className="flex flex-column">
                                            <div className="border-2 border-dashed surface-border border-round surface-ground flex-auto flex justify-content-center align-items-center font-medium p-3">
                                                <div className="w-full">
                                                    <InputNumber id="numberOfCandidates"  placeholder="Number of candidates"
                                                                 className="w-full mb-3" name="numberOfCandidates"
                                                                 value={formik.values.numberOfCandidates}
                                                                 onValueChange={formik.handleChange} showButtons min={0} max={100} mode="decimal" step={1}
                                                    />
                                                    {formik.errors.numberOfCandidates && formik.touched.numberOfCandidates &&
                                                        <Message severity={"error"} text={formik.errors.numberOfCandidates}
                                                                 className="w-full mb-1"/>}

                                                    <InputNumber id="numberOfSelection"  placeholder="Number of selection"
                                                                 className="w-full mb-3" name="numberOfSelection"
                                                                 value={formik.values.numberOfSelection}
                                                                 onValueChange={formik.handleChange} showButtons min={0} max={100} mode="decimal" step={1}
                                                    />
                                                    {formik.errors.numberOfSelection && formik.touched.numberOfSelection &&
                                                        <Message severity={"error"} text={formik.errors.numberOfSelection}
                                                                 className="w-full mb-1"/>}
                                                </div>
                                            </div>
                                            <div className="flex pt-4 justify-content-between">
                                                <Button label="Back" severity="secondary" icon="pi pi-arrow-left" type="button" raised
                                                        onClick={() => stepperRef.current.prevCallback()}/>

                                                <Button label="Register Election" icon="pi pi-send" type="submit" raised/>
                                            </div>
                                        </div>
                                    </StepperPanel>
                                </Stepper>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
    )
}