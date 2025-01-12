"use client"

import {useRef, useState} from "react";
import {Card} from "primereact/card";
import {InputText} from "primereact/inputtext";
import {classNames} from "primereact/utils";
import {Message} from "primereact/message";
import {Button} from "primereact/button";
import {useRouter} from "next/navigation";
import axios from "axios";

export default function VoterRegistrationForm() {
    const [formData, setFormData] = useState({
        idnp: '',
        fullName: ''
    })

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState(false);
    const router = useRouter();

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
        setError('');
        setSuccess(false);
    }

    const validateForm = () => {
        if (!formData.idnp || !formData.fullName) {
            setError('All fields are required');
            return false;
        }
        if (formData.idnp.length !== 13 || !/^\d+$/.test(formData.idnp)) {
            setError('IDNP must be exactly 13 digits');
            return false;
        }
        if (formData.fullName.length < 2) {
            setError('Fullname must be at least 2 characters long');
            return false;
        }
        return true;
    }

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (!validateForm()) return;

        setLoading(true);
        setError('');

        try {
            const url = `http://localhost:3000/invoke?channelid=mychannel&chaincodeid=basic&function=RegisterVoter&args=${formData.idnp}&args=${formData.fullName}`;

            const response = await axios.post(url);


            if (response.data.toString() !== 'voter registered successfully') {
                throw new Error('Registration failed');
            }

            setSuccess(true);
            setTimeout(()=> router.push('/election'), 5000);
            setFormData({ idnp: '', fullName: ''});
        } catch (err) {
            setError(err.message || 'Failed to register voter');
        } finally {
            setLoading(false);
        }
    }

    const cardHeader = (
        <div className="flex align-items-center justify-content-center">
            <h4 className="m-0 text-center">Voter Registration</h4>
        </div>
    );

    return (
        <div className="flex align-items-center justify-content-center min-h-screen">
            <Card title={cardHeader} className="w-full md:w-30rem">
                <form onSubmit={handleSubmit} className="flex flex-column gap-4">
                    <div className="flex flex-column gap-2">
                        <label htmlFor="idnp" className="font-bold">
                            IDNP
                        </label>
                        <InputText
                            id="idnp"
                            name="idnp"
                            value={formData.idnp}
                            onChange={handleChange}
                            placeholder="Enter 13-digit IDNP"
                            maxLength={13}
                            className={classNames({'p-invalid': error && !formData.idnp})}
                        />
                    </div>

                    <div className="flex flex-column gap-2">
                        <label htmlFor="fullName" className="font-bold">
                            Full Name
                        </label>
                        <InputText
                            id="fullName"
                            name="fullName"
                            value={formData.fullName}
                            onChange={handleChange}
                            placeholder="Enter full name"
                            className={classNames({'p-invalid': error && !formData.fullName})}
                        />
                    </div>

                    {error && (
                        <Message severity="error" text={error} className="w-full"/>
                    )}

                    {success && (
                        <Message severity="success" text="Voter registered successfully!" className="w-full"/>
                    )}

                    <Button
                        type="submit"
                        label={loading ? 'Registering...' : 'Register Voter'}
                        icon={loading ? 'pi pi-spin pi-spinner' : 'pi pi-user-plus'}
                        loading={loading}
                        disabled={loading}
                    />
                </form>
            </Card>
        </div>
    )
}