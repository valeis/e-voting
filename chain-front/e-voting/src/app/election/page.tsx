"use client"

import "primereact/resources/themes/lara-light-cyan/theme.css";
import {useQuery} from "@tanstack/react-query";
import axios from "axios";
import {ProgressSpinner} from "primereact/progressspinner";
import {Card} from "primereact/card";
import {Button} from "primereact/button";
import {useRouter} from "next/navigation";


interface Election {
    ID: number;
    title: string;
    description: string;
    type: string;
    startDate: string;
    endDate: string;
    numberOfAuthAttempts: string;
    numberOfCandidates: string;
    numberOfSelection: string;
    authMethod: string;
}

interface IGetElectionsResponse {
    elections: Election[];
    message: string,
}

interface IGetCandidatesResponse {
    message: string;
    candidates: string;
}

export default function ElectionPage() {
    const router = useRouter();
    const { data, isLoading, error } = useQuery({
        queryKey: ['elections'],
        queryFn: async ()  => {
            const response = await axios.get<IGetElectionsResponse>('http://localhost:3000/elections');
            return response.data;
        },
    });

    const getCandidates = async (electionId: number) => {
        const response = await axios.get<IGetCandidatesResponse>(`http://localhost:3000/elections/candidates/${electionId}`);
        return response.data.candidates;
    }

    if (isLoading) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <ProgressSpinner />
            </div>
        );
    }

    if (error) {
        return (
            <div className="text-center text-red-500 p-4">
                Error loading elections. Please try again later.
            </div>
        );
    }

    const getStatus = (startDate: string, endDate: string) => {
        const now = new Date();
        const start = new Date(startDate);
        const end = new Date(endDate);

        if (now < start) return 'Upcoming';
        if (now > end) return 'Completed';
        return 'Active'
    };

    const CandidateButton = ({ electionId}: { electionId: number }) => {
        const { data: candidates, isLoading } = useQuery({
            queryKey: ['candidates', electionId],
            queryFn: () => getCandidates(electionId),
        });

        if (isLoading) {
            return <Button label="Loading..." className="w-full" disabled/>
        }

        if (candidates && candidates.length > 0) {
            return (
                <div className="flex flex-wrap gap-2 justify-content-center">
                    <Button
                        label="View candidates"
                        type="button"
                        raised
                        onClick={() => router.push(`/election/candidates/${electionId}`)}
                    />
                    <Button
                        label="See result"
                        type="button"
                        raised
                        onClick={() => router.push(`/election/results/${electionId}`)}
                    />
                </div>
            )
        }

        return (
            <Button
                label="Register candidates"
                className="w-full"
                type="button"
                raised
                onClick={() => router.push(`/election/register-candidates/${electionId}`)}
            />
        )
    }

    return (
        <div className="p-6">
            <div className="flex justify-content-between items-center w-full mb-6">
                <h1 className="text-3xl font-bold">Elections</h1>
                <Button
                    label="Create election"
                    type="button"
                    icon="pi pi-plus-circle"
                    severity="success"
                    className="create-election-button"
                    raised
                    onClick={() => router.push(`/election/create`)}
                />
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {data?.elections.map((election: Election) => (
                    <Card
                        key={election.ID}
                        className="shadow-lg"
                        title={
                            <div className="flex justify-between items-center">
                                <span className="text-xl font-bold px-2">{election.title}</span>
                                <span className={`px-3 py-1 rounded-full text-sm ${
                                    getStatus(election.startDate, election.endDate) === 'Active'
                                        ? 'bg-green-100 text-green-800'
                                        : getStatus(election.startDate, election.endDate) === 'Upcoming'
                                            ? 'bg-blue-100 text-blue-800'
                                            : 'bg-gray-100 text-gray-800'
                                }`}>
                                    {getStatus(election.startDate, election.endDate)}
                                </span>
                            </div>
                        }
                    >
                        <div className="space-y-4">
                            <p className="text-gray-600">{election.description}</p>
                            <div className="grid grid-cols-2 gap-4 text-sm">
                                <div>
                                    <p className="font-semibold">Type</p>
                                    <p>{election.type}</p>
                                </div>
                                <div>
                                    <p className="font-semibold">Auth Method</p>
                                    <p>{election.authMethod}</p>
                                </div>
                                <div>
                                    <p className="font-semibold">Start Date</p>
                                    <p>{new Date(election.startDate).toLocaleDateString()}</p>
                                </div>
                                <div>
                                    <p className="font-semibold">End Date</p>
                                    <p>{new Date(election.endDate).toLocaleDateString()}</p>
                                </div>
                            </div>
                            <div className="pt-4">
                                <CandidateButton electionId={election.ID} />
                            </div>
                        </div>
                    </Card>
                ))}
            </div>
        </div>
    )
}