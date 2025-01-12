"use client"

import "primereact/resources/themes/lara-light-cyan/theme.css";
import {useRouter} from "next/navigation";
import {use, useEffect, useState} from "react";
import {useQuery} from "@tanstack/react-query";
import axios, {AxiosRequestHeaders} from "axios";
import {ProgressSpinner} from "primereact/progressspinner";
import {Card} from "primereact/card";
import {InputText} from "primereact/inputtext";
import {FileUpload} from "primereact/fileupload";
import {Button} from "primereact/button";
import {Accordion, AccordionTab} from "primereact/accordion";

interface Candidate {
    firstName: string;
    lastName: string;
    age: string;
    party: string;
    photo: string | null;
}

interface Election {
    ID: number;
    title: string;
    numberOfCandidates: string;
}

export default function CandidateRegistrationPage({ params }: { params: Promise<{ electionId: string }> }) {
    const router = useRouter();
    const resolvedParams = use(params);
    const [candidates, setCandidates] = useState<Candidate[]>([]);

    const { data: election, isLoading } = useQuery({
        queryKey: ['election', resolvedParams.electionId],
        queryFn: async () => {
            const response = await axios.get<{ election: Election }>(
                `http://localhost:3000/elections/${resolvedParams.electionId}`
            );
            return response.data.election;
        }
    });

    useEffect(() => {
        if (election) {
            const  initialCandidates = Array(parseInt(election.numberOfCandidates)).fill({
                firstName: '',
                lastName: '',
                age: '',
                party: '',
                photo: null,
            });
            setCandidates(initialCandidates);
        }
    }, [election]);

    const handleInputChange = (index: number, field: keyof Candidate, value: string | File) => {
        const updatedCandidates = [...candidates];
        updatedCandidates[index] = {
            ...updatedCandidates[index],
            [field]: value,
        };
        setCandidates(updatedCandidates);
    }

    const handleImageUpload = async (index: number, event: any) => {
        const file = event.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = () => {
             const base64String = reader.result as string;
             handleInputChange(index, 'photo', base64String.split(",")[1]);
            };
            reader.readAsDataURL(file);
        }
    };

    const handleSubmit = async () => {
        try {
            const payload = {
                electionId: Number(resolvedParams.electionId),
                candidates: candidates.map((candidate) => ({
                    firstName: candidate.firstName,
                    lastName: candidate.lastName,
                    age: Number(candidate.age),
                    party: candidate.party,
                    photo: candidate.photo,
                })),
            };

            await axios.post(`http://localhost:3000/elections/candidates`, payload, {
                headers: {
                    'Content-Type': 'application/json',
                } as AxiosRequestHeaders
            });
            router.push('/election');
        } catch (error) {
            console.error('Error registering candidates:', error);
        }
    };

    if (isLoading) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <ProgressSpinner />
            </div>
        )
    }

    return (
        <div className="p-6">
            <h1 className="text-3xl font-bold mb-6">
                Register Candidates for {election?.title}
            </h1>
            <div className="space-y-6">
                <Accordion>
                    {candidates.map((candidate, index) => (
                        <AccordionTab key={index} header={`Candidate ${index + 1}`}>
                            <Card key={index} className="shadow-lg">
                                <h2 className="text-xl font-semibold mb-4">Candidate {index + 1}</h2>
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                    <div className="space-y-2">
                                        <InputText
                                            value={candidate.firstName}
                                            onChange={(e) => handleInputChange(index, 'firstName', e.target.value)}
                                            placeholder="First Name"
                                            className="w-full"
                                        />
                                    </div>
                                    <div className="space-y-2">
                                        <InputText
                                            value={candidate.lastName}
                                            onChange={(e) => handleInputChange(index, 'lastName', e.target.value)}
                                            placeholder="Last Name"
                                            className="w-full"
                                        />
                                    </div>
                                    <div className="space-y-2">
                                        <InputText
                                            value={candidate.age}
                                            onChange={(e) => handleInputChange(index, 'age', e.target.value)}
                                            placeholder="Age"
                                            type="number"
                                            className="w-full"
                                        />
                                    </div>
                                    <div className="space-y-2">
                                        <InputText
                                            value={candidate.party}
                                            onChange={(e) => handleInputChange(index, 'party', e.target.value)}
                                            placeholder="Party"
                                            className="w-full"
                                        />
                                    </div>
                                    <div className="space-y-2 col-span-2">
                                        <FileUpload
                                            mode="basic"
                                            chooseLabel="Select Photo"
                                            accept="image/*"
                                            maxFileSize={1000000}
                                            onSelect={(e) => handleImageUpload(index, e)}
                                            className="w-full"
                                        />
                                    </div>
                                </div>
                            </Card>
                        </AccordionTab>
                    ))}
                </Accordion>
                <div className="flex justify-end mt-3">
                    <Button
                        label="Cancel"
                        onClick={() => router.push(`/election`)}
                        className="w-48"
                        raised
                        severity="secondary"
                    />
                    &nbsp;
                    <Button
                        label="Register All Candidates"
                        onClick={handleSubmit}
                        className="w-48"
                        raised
                        severity="success"
                    />
                </div>
            </div>
        </div>
    );
}